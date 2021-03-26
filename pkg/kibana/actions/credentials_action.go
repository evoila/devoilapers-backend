package actions

import (
	"OperatorAutomation/pkg/core/action"
	kbCommon "OperatorAutomation/pkg/kibana/common"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to get es credentials
func GetCredentialsAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Get credentials",
		UniqueCommand: "cmd_kb_get_credentials",
		Placeholder:   nil,
		ActionExecuteCallback: func(i interface{}) (interface{}, error) {
			return GetCredentials(service)
		},
	}
}

func GetCredentials(kb *kbCommon.KibanaServiceInformations) (interface{}, error) {
	// Use elastic search reference
	secretName := kb.ClusterInstance.Spec.ElasticsearchRef.Name + "-es-elastic-user"

	secret, err := kb.K8sApi.ClientSet.CoreV1().Secrets(kb.ClusterInstance.Namespace).Get(context.TODO(), secretName, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	passwordBytes := secret.Data["elastic"]
	return &action_dtos.CredentialsDto{Username: "elastic", Password: string(passwordBytes)}, nil
}
