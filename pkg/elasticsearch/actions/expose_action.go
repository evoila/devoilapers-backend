package actions

import (
	"OperatorAutomation/pkg/core/action"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	"context"
	kubernetesError "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateExposeToggleAction(service *esCommon.ElasticsearchServiceInformations) action.ToggleAction {
	return action.CreateToggleAction(
		"Exposed",
		"cmd_es_expose_toggle",
		func() (bool, error) {
			return IsExposed(service)
		},
		func() (interface{}, error) {
			return nil, Expose(service)
		},
		func() (interface{}, error) {
			return nil, Hide(service)
		})
}

// Creates an action to deliver port informations about the service
func CreateGetExposeInformationAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose infos",
		UniqueCommand: "cmd_es_get_expose_info",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return GetExposeInformation(service)
		},
	}
}

// Delivers information about the ingress
func IsExposed(es *esCommon.ElasticsearchServiceInformations) (bool, error) {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	_, err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Get(
		context.TODO(),
		ingressName,
		v1.GetOptions{})


	if kubernetesError.IsNotFound(err) {
		logger.RTrace(err, "Ingress does not exists for elasticsearch with name "+es.ClusterInstance.ObjectMeta.Name)
		return false, nil
	}

	if err != nil {
		logger.RError(err, "Could not receive ingress for "+es.ClusterInstance.ObjectMeta.Name)
		return false, err
	}

	return true, nil
}

// Delivers information about the ingress
func GetExposeInformation(es *esCommon.ElasticsearchServiceInformations) (*action_dtos.ExposeInformations, error) {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	ingress, err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Get(
		context.TODO(),
		ingressName,
		v1.GetOptions{})

	if kubernetesError.IsNotFound(err) {
		logger.RTrace(err, "Ingress does not exists for elasticsearch with name "+es.ClusterInstance.ObjectMeta.Name)
		return &action_dtos.ExposeInformations{Status: "Not exposed", Host: "Unknown"}, nil
	}

	if err != nil {
		logger.RError(err, "Could not receive ingress for "+es.ClusterInstance.ObjectMeta.Name)
		return nil, err
	}

	host := ingress.Spec.Rules[0].Host
	return &action_dtos.ExposeInformations{Status: "Exposed", Host: host}, err
}

// Reverts the expose action by removing the ingress
func Hide(es *esCommon.ElasticsearchServiceInformations) error {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Delete(context.TODO(), ingressName, v1.DeleteOptions{})

	if err != nil {
		logger.RError(err, "Could not hide ingress for "+es.ClusterInstance.ObjectMeta.Name)
	}

	return err
}

// Open a port to connect to the elasticsearch from outside
func Expose(es *esCommon.ElasticsearchServiceInformations) error {
	tlsSecretName := es.ClusterInstance.ObjectMeta.Name + "-es-http-certs-public"
	serviceName := es.ClusterInstance.ObjectMeta.Name + "-es-http"
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	hostname := es.ClusterInstance.Name + "." + es.ClusterInstance.Namespace + "." + es.Hostname

	trueValue := true
	ownerRef := v1.OwnerReference{
		UID:                es.ClusterInstance.UID,
		APIVersion:         esCommon.GroupName + "/" + esCommon.GroupVersion,
		Name:               es.ClusterInstance.Name,
		Kind:               "Elasticsearch",
		Controller:         &trueValue,
		BlockOwnerDeletion: &trueValue,
	}

	_, err := es.K8sApi.CreateIngressWithHttpsBackend(
		ingressName,
		es.ClusterInstance.Namespace,
		hostname,
		tlsSecretName,
		serviceName,
		9200,
		ownerRef)

	return err
}
