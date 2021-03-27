package actions

import (
	"OperatorAutomation/pkg/core/action"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to get es credentials
func GetCredentialsAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Get credentials",
		UniqueCommand: "cmd_es_get_credentials",
		Placeholder:   nil,
		ActionExecuteCallback: func(i interface{}) (interface{}, error) {
			return GetCredentials(service)
		},
	}
}

func GetCredentials(es *esCommon.ElasticsearchServiceInformations) (interface{}, error) {
	secretName := es.ClusterInstance.Name + "-es-elastic-user"

	secret, err := es.K8sApi.ClientSet.CoreV1().Secrets(es.ClusterInstance.Namespace).Get(context.TODO(), secretName, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	passwordBytes := secret.Data["elastic"]
	return &action_dtos.CredentialsDto{Username: "elastic", Password: string(passwordBytes)}, nil
}
