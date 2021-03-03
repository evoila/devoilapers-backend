package provider

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/utils"
	"io/ioutil"
)

type BasicProvider struct {
	Template     *service.IServiceTemplate
	Host         string
	CaPath       string
	Description  string
	Image        string
	ProviderType string
}

func CreateCommonProvider(
	host string,
	caPath string,
	templatePath string,

	providerType string,
	description string,
	image string,

) BasicProvider {
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		panic("Yaml template could not be found under path: " + templatePath)
	}

	var template service.IServiceTemplate = service.ServiceTemplate{
		Yaml:              string(templateData),
		ImportantSections: []string{"metadata.name"},
	}

	return BasicProvider{
		Template:     &template,
		Host:         host,
		CaPath:       caPath,
		Description:  description,
		Image:        image,
		ProviderType: providerType,
	}
}

func (cp BasicProvider) GetTemplate(auth common.IKubernetesAuthInformation) *service.IServiceTemplate {
	originalTemplate := *cp.Template
	yamlTemplate := originalTemplate.GetYAML()
	yamlTemplate = utils.FillWithData(auth, yamlTemplate)

	var template service.IServiceTemplate = service.ServiceTemplate{
		Yaml:              yamlTemplate,
		ImportantSections: originalTemplate.GetImportantSections(),
	}

	return &template
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
