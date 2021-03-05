package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres/actions/dtos"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func CreateGetCredentialsAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.Action{
		Name:          "Get credentials",
		UniqueCommand: "cmd_pg_get_credentials",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return GetDatabaseCredentials(service)
		},
	}
}

func GetDatabaseCredentials(pg *pgCommon.PostgresServiceInformations) (*dtos.ClusterCredentialsDto, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	// Construct the secret name using the crunchy operator convention
	nameOfSecret := pg.ClusterInstance.Name + "-" + pg.ClusterInstance.Spec.User + "-secret"
	secret, err := api.ClientSet.CoreV1().Secrets(pg.Auth.GetKubernetesNamespace()).Get(
		context.TODO(),
		nameOfSecret,
		metav1.GetOptions{},
	)

	if err != nil {
		return nil, err
	}

	// Convert port to int
	port, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return nil, err
	}

	return &dtos.ClusterCredentialsDto{
		Username: string(secret.Data["username"]),
		Password: string(secret.Data["password"]),
		Port:     port,
	}, nil
}
