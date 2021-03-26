package common

import (
	"OperatorAutomation/pkg/kubernetes"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
)

type KibanaServiceInformations struct {
	Hostname         string
	ClusterInstance  *v1.Kibana
	K8sApi			 *kubernetes.K8sApi
	CrdClient        *kubernetes.CommonCrdApi
}

