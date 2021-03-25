package actions

import (
	"OperatorAutomation/pkg/core/action"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to deliver port informations about the service
func CreateGetExposeInformationAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose infos",
		UniqueCommand: "cmd_pg_get_expose_info",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return GetExposeInformation(service)
		},
	}
}

// Creates an action to expose the service with a random port
func CreateExposeAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose",
		UniqueCommand: "cmd_pg_expose",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, Expose(service)
		},
	}
}

// Creates an action to remove the exposure
func DeleteExposeAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Hide",
		UniqueCommand: "cmd_pg_hide",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, Hide(service)
		},
	}
}

// Delivers information about the ingress
func GetExposeInformation(es *esCommon.ElasticsearchServiceInformations) (*action_dtos.ExposeInformations, error) {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	ingress, err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Get(
		context.TODO(),
		ingressName,
		v1.GetOptions{})

	if err != nil {
		logger.RError(err, "Could not receive ingress for " + es.ClusterInstance.ObjectMeta.Name)
		return nil, err
	}

	host := ingress.Spec.Rules[0].Host
	return &action_dtos.ExposeInformations{Host: host}, err
}

// Reverts the expose action by removing the ingress
func Hide(es *esCommon.ElasticsearchServiceInformations) error {

	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Delete(context.TODO(), ingressName, v1.DeleteOptions{})

	if err != nil {
		logger.RError(err, "Could not hide ingress for " + es.ClusterInstance.ObjectMeta.Name)
	}

	return err
}

// Open a port to connect to the elasticsearch from outside
func Expose(es *esCommon.ElasticsearchServiceInformations) error {

	tlsSecretName := es.ClusterInstance.ObjectMeta.Name + "-es-http-certs-public"
	serviceName := es.ClusterInstance.ObjectMeta.Name + "-es-http"
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-es-ingress"
	hostname := es.ClusterInstance.Name + "." + es.ClusterInstance.Namespace + "." + es.Hostname

	_, err := es.K8sApi.CreateIngressWithHttpsBackend(ingressName, es.ClusterInstance.Namespace, hostname, tlsSecretName, serviceName, 9200)

	return err
}
