package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	"fmt"

	"OperatorAutomation/pkg/kubernetes"
	"strings"
	"testing"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	"github.com/stretchr/testify/assert"
)

func CreateElasticSearchProvider(t *testing.T) (*elasticsearch.ElasticsearchProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

	return &esProvider, config
}

func CreateEsDummyService(t *testing.T, provider *elasticsearch.ElasticsearchProvider,
	user common.IKubernetesAuthInformation) (*elasticsearch.ElasticSearchService, [2]string) {
	yaml := (*provider.GetTemplate(user)).GetYAML()
	er := provider.CreateService(user, yaml)
	if er != nil {
		t.Error("Fail to create ES instance")
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

	service, _ := (*iService).(elasticsearch.ElasticSearchService)
	return &service, infos
}

func Test_Elastic_Expose(t *testing.T) {
	provider, config := CreateElasticSearchProvider(t)
	user := config.Users[0]
	k8sapi, _ := kubernetes.GenerateK8sApiFromToken(provider.Host, provider.CaPath, user.GetKubernetesAccessToken())
	service, infos := CreateEsDummyService(t, provider, user)
	servicename := infos[0] + "-es-http"
	ingressname := "my-test-ingress"
	exposeinfo := dtos.ExposeInformation{IngressName: ingressname}

	service.ExecuteExposeAction(&exposeinfo)
	ingress, _ := k8sapi.GetIngress(user.GetKubernetesNamespace(), ingressname)
	assert.True(t, k8sapi.ExistingServiceInIngress(ingress, servicename), "ingress not found")

	service.ExecuteUnexposeAction(&exposeinfo)
	ingress, _ = k8sapi.GetIngress(user.GetKubernetesNamespace(), ingressname)
	assert.False(t, k8sapi.ExistingServiceInIngress(ingress, servicename), "ingress should have been removed")

	k8sapi.DeleteServiceFromIngress(user.KubernetesNamespace, ingressname, servicename)
	provider.DeleteService(user, infos[0])
}

func Test_Elastic_Rescale(t *testing.T) {
	provider, config := CreateElasticSearchProvider(t)
	user := config.Users[0]
	k8sapi, _ := kubernetes.GenerateK8sApiFromToken(provider.Host, provider.CaPath, user.GetKubernetesAccessToken())

	service, infos := CreateEsDummyService(t, provider, user)
	statefulsetname := infos[0] + "-es-" + infos[1]
	t.Log("sts name", statefulsetname)
	scalenum := int32(2)

	WaitForEsTobeReady(provider, user, infos[0])

	scaleinfo := dtos.ScaleInformation{ReplicasCount: int32(scalenum)}
	service.ExecuteRescaleAction(&scaleinfo)

	scale, _ := k8sapi.GetStatefulSetScale(user.KubernetesNamespace, statefulsetname)
	//t.Log("Es instance has been created and scaled", scale)
	assert.Equal(t, scalenum, scale.Spec.Replicas, "should be "+fmt.Sprintf("%d", scalenum)+" but found "+string(scale.Spec.Replicas))

	provider.DeleteService(user, infos[0])
}

func WaitForEsTobeReady(provider *elasticsearch.ElasticsearchProvider, user common.IKubernetesAuthInformation, esname string) {
	esCrdAPI, _ := kubernetes.CreateCommonCrdApi(provider.Host, provider.CaPath, user.GetKubernetesAccessToken(), elasticsearch.GroupName, elasticsearch.GroupVersion)
	instance := v1.Elasticsearch{}
	for {
		esCrdAPI.Get(user.GetKubernetesNamespace(), esname, elasticsearch.RessourceName, &instance)
		if instance.Status.Health == v1.ElasticsearchGreenHealth {
			break
		}
	}
}

/*func Test_Elastic_Other(t *testing.T) {
	t.Log("Start other test")
	provider, config := CreateElasticSearchProvider(t)
	user := config.Users[0]
	api, _ := kubernetes.GenerateK8sApiFromToken(provider.Host, provider.CaPath, user.GetKubernetesAccessToken())
	ingressName := "test-ingress"
	serviceName := "elasticsearch-sample-es-http"
	namespace := user.KubernetesNamespace
	ingress, err := api.GetIngress(namespace, ingressName)
	if err != nil {
	}
	l := len(ingress.Spec.Rules[0].HTTP.Paths)
	t.Log("ingress:", l)
	if l > 1 {
		t.Log("before:", ingress.Spec.Rules[0].HTTP.Paths)
		for i, path := range ingress.Spec.Rules[0].HTTP.Paths {
			if path.Backend.ServiceName == serviceName {
				t.Log("foundsth")
				ingress.Spec.Rules[0].HTTP.Paths[i] = ingress.Spec.Rules[0].HTTP.Paths[l-1]
				ingress.Spec.Rules[0].HTTP.Paths = ingress.Spec.Rules[0].HTTP.Paths[:l-1]
				break
			}
		}
		t.Log("after:", ingress.Spec.Rules[0].HTTP.Paths)
		_, err = api.V1beta1Client.Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	} else if l == 1 {
		if ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName == serviceName {
			err = api.V1beta1Client.Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
		}
	}
}*/
