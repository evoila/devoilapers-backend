package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kubernetes"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	dtos2 "OperatorAutomation/pkg/postgres/dtos"
	"errors"
	pgV1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"strings"
)

const PgreplicaRessource = "pgreplicas"

func CreateScaleAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Scale",
		UniqueCommand: "cmd_pg_scale",
		Placeholder:   CreatePlaceholder(service),
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return scaleCluster(service, placeholder.(*dtos2.ClusterScaleDto))
		},
	}
}

// Generate placeholder which represents the required input fields
func CreatePlaceholder(pg *pgCommon.PostgresServiceInformations) *dtos2.ClusterScaleDto {
	replicas, err := getReplicas(pg)
	numberOfReplicas := 0

	if err == nil {
		numberOfReplicas = len(replicas.Items)
	}

	return &dtos2.ClusterScaleDto{
		NumberOfReplicas: numberOfReplicas,
	}
}

func getReplicas(pg *pgCommon.PostgresServiceInformations) (*pgV1.PgreplicaList, error) {
	postgresReplicas := pgV1.PgreplicaList{}

	clusterSelector := metav1.ListOptions{
		LabelSelector: "pg-cluster=" + pg.ClusterInstance.Name,
	}

	err := pg.CrdClient.ListWithOptions(
		pg.Auth.GetKubernetesNamespace(),
		PgreplicaRessource,
		&clusterSelector,
		&postgresReplicas)

	return &postgresReplicas, err
}

func scaleUp(pg *pgCommon.PostgresServiceInformations, numberOfNewReplicas int) error {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, pg.Auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	pgouser, found := pg.ClusterInstance.Labels["pgouser"]
	if !found {
		return errors.New("Invalid postgres cluster. Label \"pgouser\" is missing.")
	}

	for i := 0; i < numberOfNewReplicas; i++ {
		newUuid := strings.ToLower(strings.Replace(uuid.New().String(), "-", "", -1))
		name := pg.ClusterInstance.Name + "-replica-" + newUuid
		replica := pgV1.Pgreplica{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pgreplica",
				APIVersion: pg.ClusterInstance.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name":       name,
					"pg-cluster": pg.ClusterInstance.Spec.ClusterName,
					"pgouser":    pgouser,
				},
				Name:      name,
				Namespace: pg.ClusterInstance.Namespace,
			},
			Spec: pgV1.PgreplicaSpec{
				ClusterName: pg.ClusterInstance.Spec.ClusterName,
				Name:        name,
				Namespace:   pg.ClusterInstance.Namespace,
				ReplicaStorage: pgV1.PgStorageSpec{
					AccessMode:         pg.ClusterInstance.Spec.ReplicaStorage.AccessMode,
					MatchLabels:        pg.ClusterInstance.Spec.ReplicaStorage.MatchLabels,
					Name:               name,
					Size:               pg.ClusterInstance.Spec.ReplicaStorage.Size,
					StorageClass:       pg.ClusterInstance.Spec.ReplicaStorage.StorageClass,
					StorageType:        pg.ClusterInstance.Spec.ReplicaStorage.StorageType,
					SupplementalGroups: pg.ClusterInstance.Spec.ReplicaStorage.SupplementalGroups,
				},
				UserLabels: pg.ClusterInstance.Spec.UserLabels,
			},
		}

		yamlData, err := yaml.Marshal(&replica)
		if err != nil {
			return errors.New("Replica could not serialized to yaml.")
		}

		_, err = api.Apply(yamlData)
		if err != nil {
			return err
		}
	}

	return nil
}

func scaleDown(
	pg *pgCommon.PostgresServiceInformations,
	numberOfNewReplicasToDelete int,
	replicas *pgV1.PgreplicaList) error {

	for i := 0; i < numberOfNewReplicasToDelete; i++ {
		replica := replicas.Items[i]
		err := pg.CrdClient.Delete(replica.Namespace, replica.Name, PgreplicaRessource)
		if err != nil {
			return err
		}
	}

	return nil
}

func scaleCluster(pg *pgCommon.PostgresServiceInformations, dto *dtos2.ClusterScaleDto) (interface{}, error) {
	if dto.NumberOfReplicas < 0 {
		return nil, errors.New("invalid total number of replicas")
	}

	currentReplicas, err := getReplicas(pg)
	if err != nil {
		return nil, err
	}

	if len(currentReplicas.Items) == dto.NumberOfReplicas {
		return nil, errors.New("neither number of replicas nor size has changed")
	} else if len(currentReplicas.Items) > dto.NumberOfReplicas {
		return nil, scaleDown(pg, len(currentReplicas.Items)-dto.NumberOfReplicas, currentReplicas)
	} else {
		return nil, scaleUp(pg, dto.NumberOfReplicas-len(currentReplicas.Items))
	}
}
