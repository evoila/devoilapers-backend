package actions

import (
	"OperatorAutomation/pkg/core/action"
	kbCommon "OperatorAutomation/pkg/kibana/common"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	"context"
	kubernetesError "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)



func CreateExposeToggleAction(service *kbCommon.KibanaServiceInformations) action.ToggleAction {
	return action.CreateToggleAction(
		"Exposed",
		"cmd_kb_expose_toggle",
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

// Delivers information about the ingress
func IsExposed(kb *kbCommon.KibanaServiceInformations) (bool, error) {
	ingressName := kb.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	_, err := kb.K8sApi.V1beta1Client.Ingresses(kb.ClusterInstance.Namespace).Get(
		context.TODO(),
		ingressName,
		v1.GetOptions{})


	if kubernetesError.IsNotFound(err) {
		logger.RTrace(err, "Ingress does not exists for kibana with name "+kb.ClusterInstance.ObjectMeta.Name)
		return false, nil
	}

	if err != nil {
		logger.RError(err, "Could not receive ingress for "+ kb.ClusterInstance.ObjectMeta.Name)
		return false, err
	}

	return true, nil
}


// Delivers information about the ingress
func GetExposeInformation(kb *kbCommon.KibanaServiceInformations) (*action_dtos.ExposeInformations, error) {
	ingressName := kb.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	ingress, err := kb.K8sApi.V1beta1Client.Ingresses(kb.ClusterInstance.Namespace).Get(
		context.TODO(),
		ingressName,
		v1.GetOptions{})

	if kubernetesError.IsNotFound(err) {
		logger.RTrace(err, "Ingress does not exists for kibana with name "+ kb.ClusterInstance.ObjectMeta.Name)
		return &action_dtos.ExposeInformations{Status: "Not exposed", Host: "Unknown"}, nil
	}

	if err != nil {
		logger.RError(err, "Could not receive ingress for " + kb.ClusterInstance.ObjectMeta.Name)
		return nil, err
	}

	host := ingress.Spec.Rules[0].Host
	return &action_dtos.ExposeInformations{Status: "Exposed", Host: host}, err
}

// Reverts the expose action by removing the ingress
func Hide(kb *kbCommon.KibanaServiceInformations) error {
	ingressName := kb.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	err := kb.K8sApi.V1beta1Client.Ingresses(kb.ClusterInstance.Namespace).Delete(context.TODO(), ingressName, v1.DeleteOptions{})

	if err != nil {
		logger.RError(err, "Could not hide ingress for " + kb.ClusterInstance.ObjectMeta.Name)
	}

	return err
}

// Open a port to connect to the elasticsearch from outside
func Expose(kb *kbCommon.KibanaServiceInformations) error {
	tlsSecretName := kb.ClusterInstance.ObjectMeta.Name + "-kb-http-certs-public"
	serviceName := kb.ClusterInstance.ObjectMeta.Name + "-kb-http"
	ingressName := kb.ClusterInstance.ObjectMeta.Name + "-kb-ingress"
	hostname := kb.ClusterInstance.Name + "." + kb.ClusterInstance.Namespace + "." + kb.Hostname

	trueValue := true
	ownerRef := v1.OwnerReference{
		UID: kb.ClusterInstance.UID,
		APIVersion:         kbCommon.GroupName+"/"+kbCommon.GroupVersion,
		Name:               kb.ClusterInstance.Name,
		Kind:             	"Kibana",
		Controller:         &trueValue,
		BlockOwnerDeletion: &trueValue,
	}

	_, err := kb.K8sApi.CreateIngressWithHttpsBackend(
		ingressName,
		kb.ClusterInstance.Namespace,
		hostname,
		tlsSecretName,
		serviceName,
		5601,
		ownerRef)

	return err
}
