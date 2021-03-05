package common

import (
	"OperatorAutomation/pkg/core/common"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

type PostgresServiceInformations struct {
	Host   string
	CaPath string
	Auth   common.IKubernetesAuthInformation
	ClusterInstance *v1.Pgcluster
}