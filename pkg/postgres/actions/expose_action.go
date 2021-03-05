package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres/actions/dtos"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"context"
	"errors"
	coreApiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func CreateExposeAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.Action{
		Name:          "Exposure",
		UniqueCommand: "cmd_pg_expose",
		Placeholder:   &dtos.ClusterCredentialsDto{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return Expose(service)
		},
	}
}

func  GetService(pg *pgCommon.PostgresServiceInformations) (*kubernetes.K8sApi, *coreApiV1.Service, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())

	if err != nil {
		return api, nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: "pg-cluster=hoonah-emma",
	}

	services, err := api.ClientSet.CoreV1().Services(pg.Auth.GetKubernetesNamespace()).List(context.TODO(), opts)
	if err != nil {
		return api, nil, err
	}

	postgresServicePort, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return api, nil, err
	}

	for _, service := range services.Items {
		for _, port := range service.Spec.Ports {
			if int(port.TargetPort.IntVal) == postgresServicePort {
				return api, &service, nil
			}
		}
	}

	return api, nil, errors.New("Kubernetes service with matching name and port for postgres cluster could not be found.")
}

func Expose(pg *pgCommon.PostgresServiceInformations) (interface{}, error) {
	api, service, err := GetService(pg)
	if err != nil {
		return nil, err
	}

	_ = api
	_ = service
	return nil, nil
}
