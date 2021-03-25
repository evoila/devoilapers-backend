package common

import (
	"OperatorAutomation/pkg/kubernetes"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

type ElasticsearchServiceInformations struct {
	Hostname         string
	ClusterInstance  *v1.Elasticsearch
	K8sApi			 *kubernetes.K8sApi
	CrdClient        *kubernetes.CommonCrdApi
}

