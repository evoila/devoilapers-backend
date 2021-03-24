package common

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres/pgo"
	PostgresOperatorV1 "github.com/crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

type PostgresServiceInformations struct {
	Hostname 		 string
	HostWithPort     string
	CaPath           string
	Auth             common.IKubernetesAuthInformation
	ClusterInstance  *PostgresOperatorV1.Pgcluster
	CrdClient        *kubernetes.CommonCrdApi
	ClusterReplica   []*PostgresOperatorV1.Pgreplica
	NginxInformation kubernetes.NginxInformation
	PgoApi           *pgo.PgoApi
}
