package actions

import (
	"OperatorAutomation/pkg/core/action"
	kbCommon "OperatorAutomation/pkg/kibana/common"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to deliver port informations about the service
func CreateGetExposeInformationAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose infos",
		UniqueCommand: "cmd_kb_get_expose_info",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return GetExposeInformation(service)
		},
	}
}

// Creates an action to expose the service with a random port
func CreateExposeAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose",
		UniqueCommand: "cmd_kb_expose",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, Expose(service)
		},
	}
}

// Creates an action to remove the exposure
func DeleteExposeAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Hide",
		UniqueCommand: "cmd_kb_hide",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, Hide(service)
		},
	}
}

// Delivers information about the ingress
func GetExposeInformation(es *kbCommon.KibanaServiceInformations) (*action_dtos.ExposeInformations, error) {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
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
func Hide(es *kbCommon.KibanaServiceInformations) error {
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	err := es.K8sApi.V1beta1Client.Ingresses(es.ClusterInstance.Namespace).Delete(context.TODO(), ingressName, v1.DeleteOptions{})

	if err != nil {
		logger.RError(err, "Could not hide ingress for " + es.ClusterInstance.ObjectMeta.Name)
	}

	return err
}

// Open a port to connect to the elasticsearch from outside
func Expose(es *kbCommon.KibanaServiceInformations) error {
	tlsSecretName := es.ClusterInstance.ObjectMeta.Name + "-kb-http-certs-public"
	serviceName := es.ClusterInstance.ObjectMeta.Name + "-kb-http"
	ingressName := es.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	hostname := es.ClusterInstance.Name + "." + es.ClusterInstance.Namespace + "." + es.Hostname

	trueValue := true
	ownerRef := v1.OwnerReference{
		UID: es.ClusterInstance.UID,
		APIVersion:         kbCommon.GroupName+"/"+kbCommon.GroupVersion,
		Name:               es.ClusterInstance.Name,
		Kind:             	"Kibana",
		Controller:         &trueValue,
		BlockOwnerDeletion: &trueValue,
	}

	_, err := es.K8sApi.CreateIngressWithHttpsBackend(
		ingressName,
		es.ClusterInstance.Namespace,
		hostname,
		tlsSecretName,
		serviceName,
		5601,
		ownerRef)

	return err
}
