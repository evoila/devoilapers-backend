package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kubernetes"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos/action_dtos"
	"context"
	"errors"
	coreApiV1 "k8s.io/api/core/v1"
	kubernetesError "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func CreateExposeToggleAction(service *pgCommon.PostgresServiceInformations) action.ToggleAction {
	return action.CreateToggleAction(
		"Exposed",
		"cmd_pg_expose_toggle",
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

func IsExposed(service *pgCommon.PostgresServiceInformations) (bool, error) {
	exposeInfos, err := GetExposeInformation(service)

	if err != nil {
		return false, err
	}

	return exposeInfos.Port > 0, nil
}

// Delivers information about the exposed port
func GetExposeInformation(pg *pgCommon.PostgresServiceInformations) (*action_dtos.ClusterExposeResponseDto, error) {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.HostWithPort, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
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

	if err != nil && (kubernetesError.IsNotFound(err) || err.Error() == "no port found. The service is not exposed") {
		return &action_dtos.ClusterExposeResponseDto{Status: "Not exposed", Host: "Unknown"}, nil
	}

	if err != nil {
		return nil, err
	}

	return &action_dtos.ClusterExposeResponseDto{
		Status: "Exposed",
		Port: exposedPort,
		Host: pg.Hostname,
	}, err
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

	return nil, errors.New("kubernetes service with matching name and port for postgres cluster could not be found")
}

// Reverts the expose action by removing the port
func Hide(pg *pgCommon.PostgresServiceInformations) error {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.HostWithPort, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
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
func Expose(pg *pgCommon.PostgresServiceInformations) error {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.HostWithPort, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return  err
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

	_, err = kubernetes.NginxOpenRandomTcpPort(
		api,
		pg.NginxInformation,
		pg.ClusterInstance.Namespace,
		pg.ClusterInstance.Name,
		internalPort,
	)

	return err
}
