package provider

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"io/ioutil"
)

type BasicProvider struct {
	FormTemplate string
	YamlTemplate string
	Host         string
	CaPath       string
	Description  string
	Image        string
	ProviderType string
}

func (cp BasicProvider) OnCoreInitialized(provider []*provider.IServiceProvider) {

}

func (cp BasicProvider) GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
	panic("implement me")
}

func (cp BasicProvider) GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error) {
	panic("implement me")
}

func (cp BasicProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	panic("implement me")
}

func (cp BasicProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	panic("implement me")
}

func (cp BasicProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	panic("implement me")
}

func (cp BasicProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	panic("implement me")
}

func CreateCommonProvider(
	host string,
	caPath string,
	yamlTemplatePath string,
	formTemplatePath string,
	providerType string,
	description string,
	image string,

) BasicProvider {
	templateData, err := ioutil.ReadFile(yamlTemplatePath)
	if err != nil {
		panic("Yaml template could not be found under path: " + yamlTemplatePath)
	}

	formTemplateData, err := ioutil.ReadFile(formTemplatePath)
	if err != nil {
		panic("Form template could not be found under path: " + formTemplatePath)
	}

	return BasicProvider{
		YamlTemplate: string(templateData),
		FormTemplate: string(formTemplateData),
		Host:         host,
		CaPath:       caPath,
		Description:  description,
		Image:        image,
		ProviderType: providerType,
	}
}

func (cp BasicProvider) GetServiceDescription() string {
	return cp.Description
}

func (cp BasicProvider) GetServiceImage() string {
	return cp.Image
}

func (cp BasicProvider) GetServiceType() string {
	return cp.ProviderType
}
