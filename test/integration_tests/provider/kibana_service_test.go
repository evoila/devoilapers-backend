package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/kibana"
	"OperatorAutomation/pkg/kibana/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	"fmt"

	"OperatorAutomation/pkg/kubernetes"
	"strings"
	"testing"

	commonV1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	"github.com/stretchr/testify/assert"
)

func CreateKibanaProvider(t *testing.T) (*kibana.KibanaProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	kbProvider := kibana.CreateKibanaProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

	return &kbProvider, config
}

func CreateKbDummyService(t *testing.T, provider *kibana.KibanaProvider, user common.IKubernetesAuthInformation) (*kibana.KibanaService, [2]string,
	*elasticsearch.ElasticsearchProvider, string) {
	esProvider, _ := CreateElasticSearchProvider(t)
	esService, _ := CreateEsDummyService(t, esProvider, user)

	WaitForEsTobeReady(esProvider, user, esService.Name)

	yaml := (*provider.GetTemplate(user)).GetYAML()
	yaml = strings.Replace(yaml, "YOUR_ELASTICSEARCH_INSTANCE_NAME", esService.GetName(), 1)
	er := provider.CreateService(user, yaml)
	if er != nil {
		t.Error("Fail to create KB instance")
	}
	substrings := strings.Fields(yaml)

	infos := [2]string{"", ""}
	k := 0
	for i, s := range substrings {
		if s == "name:" || s == "count:" {
			infos[k] = substrings[(i + 1)]
			k++
			if k >= 2 {
				break
			}
		}

	}
	var iService *service.IService
	//testser, _ := esProvider.GetService(user, "derbycenter-mia")
	//t.Error("test get service", testser)
	for {
		iService, _ = provider.GetService(user, infos[0])
		if iService != nil {
			break
		}
	}
	service, _ := (*iService).(kibana.KibanaService)
	return &service, infos, esProvider, esService.Name
}

func Test_Kibana_Expose(t *testing.T) {
	provider, config := CreateKibanaProvider(t)
	user := config.Users[0]
	k8sapi, _ := kubernetes.GenerateK8sApiFromToken(provider.Host, provider.CaPath, user.GetKubernetesAccessToken())
	service, infos, esProvider, esname := CreateKbDummyService(t, provider, user)
	servicename := infos[0] + "-kb-http"
	ingressname := "my-test-ingress"
	exposeinfo := dtos.ExposeInformation{IngressName: ingressname}

	url, _ := service.ExecuteExposeAction(&exposeinfo)
	t.Log("exposed with url", url)
	ingress, _ := k8sapi.GetIngress(user.GetKubernetesNamespace(), ingressname)
	t.Log("ingress instance has been created ", ingress)
	assert.True(t, k8sapi.ExistingServiceInIngress(ingress, servicename), "ingress not found")

	service.ExecuteUnexposeAction()
	ingress, _ = k8sapi.GetIngress(user.GetKubernetesNamespace(), ingressname)
	assert.False(t, k8sapi.ExistingServiceInIngress(ingress, servicename), "ingress should have been removed")

	k8sapi.DeleteServiceFromIngress(user.KubernetesNamespace, ingressname, servicename)
	provider.DeleteService(user, infos[0])
	esProvider.DeleteService(user, esname)
}

func Test_Kibana_Rescale(t *testing.T) {
	provider, config := CreateKibanaProvider(t)
	user := config.Users[0]
	k8sapi, _ := kubernetes.GenerateK8sApiFromToken(provider.Host, provider.CaPath, user.GetKubernetesAccessToken())
	kbCrdAPI, _ := kubernetes.CreateCommonCrdApi(provider.Host, provider.CaPath, user.GetKubernetesAccessToken(), kibana.GroupName, kibana.GroupVersion)

	service, infos, esProvider, esname := CreateKbDummyService(t, provider, user)
	deploymentname := infos[0] + "-kb"
	t.Log("deployment name", deploymentname)
	scalenum := int32(2)

	instance := v1.Kibana{}
	for {
		kbCrdAPI.Get(user.KubernetesNamespace, infos[0], kibana.ResourceName, &instance)
		if instance.Status.Health == commonV1.GreenHealth {
			break
		}
	}
	scaleinfo := dtos.ScaleInformation{ReplicasCount: int32(scalenum)}
	service.ExecuteRescaleAction(&scaleinfo)

	scale, _ := k8sapi.GetDeploymentScale(user.KubernetesNamespace, deploymentname)
	t.Log("Kb instance has been created and scaled", scale)
	assert.Equal(t, scalenum, scale.Spec.Replicas, "should be "+fmt.Sprintf("%d", scalenum)+" but found "+string(scale.Spec.Replicas))

	provider.DeleteService(user, infos[0])
	esProvider.DeleteService(user, esname)
}
