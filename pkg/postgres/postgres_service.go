package postgres

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres/dtos"
	"OperatorAutomation/pkg/utils/provider"
	"context"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

type PostgresService struct {
	clusterInstance *v1.Pgcluster
	provider.BasicService

	host   string
	caPath string
	auth   common.IKubernetesAuthInformation
}

// Returns the action groups. Part of service.IService interface
func (pg PostgresService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Informations",
			Actions: []action.IAction{
				action.Action{
					Name:          "Get credentials",
					Placeholder:   &dtos.ClusterCredentialsDto{},
					UniqueCommand: "cmd_get_credentials",
					ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
						return pg.GetDatabaseCredentials()
					},
				},
			},
		},
	}
}

func (pg PostgresService) GetDatabaseCredentials() (*dtos.ClusterCredentialsDto, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.host, pg.caPath, pg.auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	nameOfSecret := pg.clusterInstance.Name + "-" + pg.clusterInstance.Spec.User + "-secret"
	secret, err := api.ClientSet.CoreV1().Secrets(pg.auth.GetKubernetesNamespace()).Get(
		context.TODO(),
		nameOfSecret,
		metav1.GetOptions{},
	)

	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(pg.clusterInstance.Spec.Port)
	if err != nil {
		return nil, err
	}

	return &dtos.ClusterCredentialsDto{
		Username: string(secret.Data["username"]),
		Password: string(secret.Data["password"]),
		Port: port,
	},  nil
}

func (pg PostgresService) GetStatus() int {

	status := pg.clusterInstance.Status.State
	if status == v1.PgclusterStateProcessed ||
		status == v1.PgclusterStateBootstrapping {
		return service.ServiceStatusPending
	} else if status == v1.PgclusterStateInitialized ||
		status == v1.PgclusterStateCreated ||
		status == v1.PgclusterStateBootstrapped {
		return service.ServiceStatusOk
	} else if status == v1.PgclusterStateBootstrapping {

	}

	return service.ServiceStatusError
}
