package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kubernetes"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	dtos2 "OperatorAutomation/pkg/postgres/dtos"
	"context"
	"errors"
	coreApiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// Creates an action to deliver port informations about the service
func CreateGetExposeInformationAction(service *pgCommon.PostgresServiceInformations) action.IAction {
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
func CreateExposeAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Expose",
		UniqueCommand: "cmd_pg_expose",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return Expose(service)
		},
	}
}

// Creates an action to remove the exposure
func DeleteExposeAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Hide",
		UniqueCommand: "cmd_pg_hide",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, Hide(service)
		},
	}
}

// Delivers information about the exposed port
func GetExposeInformation(pg *pgCommon.PostgresServiceInformations) (*dtos2.ClusterExposeResponseDto, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	internalPort, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return nil, err
	}

	exposedPort, err := kubernetes.NginxGetExposedTcpPort(
		api,
		pg.NginxInformation,
		pg.ClusterInstance.Namespace,
		pg.ClusterInstance.Name,
		internalPort,
	)

	return &dtos2.ClusterExposeResponseDto{Port: exposedPort}, err
}

// Helper method to get the service of a postgres cluster
func getService(pg *pgCommon.PostgresServiceInformations, api *kubernetes.K8sApi) (*coreApiV1.Service, error) {

	opts := metav1.ListOptions{
		LabelSelector: "pg-cluster=" + pg.ClusterInstance.Name,
	}

	services, err := api.ClientSet.CoreV1().Services(pg.Auth.GetKubernetesNamespace()).List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	postgresServicePort, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return nil, err
	}

	for _, service := range services.Items {
		for _, port := range service.Spec.Ports {
			if int(port.TargetPort.IntVal) == postgresServicePort {
				return &service, nil
			}
		}
	}

	return nil, errors.New("Kubernetes service with matching name and port for postgres cluster could not be found.")
}

// Reverts the expose action by removing the port
func Hide(pg *pgCommon.PostgresServiceInformations) error {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	// Ensure a service with matching port exists
	_, err = getService(pg, api)
	if err != nil {
		return err
	}

	// Get the cluster internal service port
	internalPort, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return err
	}

	// Finally try close the port
	return kubernetes.NginxCloseTcpPort(
		api,
		pg.NginxInformation,
		pg.ClusterInstance.Namespace,
		pg.ClusterInstance.Name,
		internalPort,
	)
}

// Open a port to connect to the db from outside
func Expose(pg *pgCommon.PostgresServiceInformations) (*dtos2.ClusterExposeResponseDto, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	// Ensure a service with matching port exists
	_, err = getService(pg, api)
	if err != nil {
		return nil, err
	}

	// Get the cluster internal service port
	internalPort, err := strconv.Atoi(pg.ClusterInstance.Spec.Port)
	if err != nil {
		return nil, err
	}

	exposedPort, err := kubernetes.NginxOpenRandomTcpPort(
		api,
		pg.NginxInformation,
		pg.ClusterInstance.Namespace,
		pg.ClusterInstance.Name,
		internalPort,
	)

	return &dtos2.ClusterExposeResponseDto{
		Port: exposedPort,
	}, err
}
