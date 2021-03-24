package controller

import (
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/pkg/utils/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type ServiceStoreController struct {
	BaseController
}


// Service store overview godoc
// @Summary Lists all possible deployable services
// @Description Lists all possible deployable services
//
// @tags Servicestore
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Success 200 {object} dtos.ServiceStoreOverviewDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /servicestore/info [get]
func (controller ServiceStoreController) HandleGetServiceStoreOverview(ctx *gin.Context) {
	logger.RTrace("Received get request to get an overview of ther service store provider")

	serviceStoreOverviewData := dtos.ServiceStoreOverviewDto{ServiceStoreItems: []dtos.ServiceStoreItemDto{}}

	for _, provider := range controller.Core.Providers {

		serviceStoreItem := dtos.ServiceStoreItemDto{
			Type:        (*provider).GetServiceType(),
			Description: (*provider).GetServiceDescription(),
			ImageSource: (*provider).GetServiceImage(),
		}

		logger.RTrace("Found provider with type " + serviceStoreItem.Type)
		serviceStoreOverviewData.ServiceStoreItems = append(serviceStoreOverviewData.ServiceStoreItems, serviceStoreItem)
	}

	logger.RTrace("Progressed all providers")
	ctx.JSON(http.StatusOK, serviceStoreOverviewData)
}

// Default Service Yaml-Template godoc
// @Summary Get the yaml for a service
// @Description  Get the yaml for a service based on the filled form and the user data
//
// @tags Servicestore
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param formresult body string true "Form-Result"
// @Param servicetype path string true "Type of service"
//
// @Success 200 {object} dtos.ServiceStoreItemYamlDto
// @Failure 400 {object} dtos.HTTPErrorDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /servicestore/yaml/{servicetype} [post]
func (controller ServiceStoreController) HandlePostServiceStoreItemYaml(ctx *gin.Context) {
	logger.RTrace("Received post request to get yaml of an service store provider")
	serviceType := ctx.Param("servicetype")

	logger.RTrace("Get provider of type " + serviceType)
	provider, err := controller.Core.GetProviderByName(serviceType)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Parse received data
	filledFormData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	// Generate yaml
	yamlObject, err := (*provider).GetYamlTemplate(userInfos, filledFormData)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Check if the type is a array. Then we separate each interface{} type in order to get
	// a valid multi-document yaml. Otherwise we just serialize it.
	var yamlBytes []byte
	switch concreteObject := yamlObject.(type) {
		case []interface{}:
			logger.RTrace("Yaml contains multiple documents. Going to marshal it separatly")

			for _, innerInterface := range concreteObject {
				// Append yaml separator
				if len(yamlBytes) > 0 {
					yamlBytes = append(yamlBytes, []byte("---\n")...)
				}

				// Convert interface to yaml
				yamlSectionBytes, err := yaml.Marshal(innerInterface)
				if err != nil {
					utils.NewError(ctx, http.StatusInternalServerError, err)
					return
				}

				// Append to the yaml
				yamlBytes = append(yamlBytes, yamlSectionBytes...)
			}
		default:
			logger.RTrace("Yaml does not contain multiple documents. Going to marshal the whole interface")
			yamlBytes, err = yaml.Marshal(yamlObject)

			if err != nil {
				utils.NewError(ctx, http.StatusInternalServerError, err)
				return
			}
	}

	logger.RTrace("Marshaling to yaml done")

	// Build response
	serviceYaml := dtos.ServiceStoreItemYamlDto{
		Yaml: string(yamlBytes),
	}

	ctx.JSON(http.StatusOK, serviceYaml)
}

// Get json form for a provider godoc
// @Summary Get the json form for a service-template
// @Description Get the default yaml file for a service-template
//
// @tags Servicestore
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param servicetype path string true "Type of service"
//
// @Success 200 {object} dtos.ServiceStoreItemFormDto
// @Failure 400 {object} dtos.HTTPErrorDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /servicestore/form/{servicetype} [get]
func (controller ServiceStoreController) HandleGetServiceStoreItemForm(ctx *gin.Context) {
	logger.RTrace("Received get request to get the form of an service store provider")
	serviceType := ctx.Param("servicetype")

	logger.RTrace("Get provider of type " + serviceType)
	provider, err := controller.Core.GetProviderByName(serviceType)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	// Generate form data
	logger.RTrace("Get json form of provider type " + serviceType)
	formData, err := (*provider).GetJsonForm(userInfos)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	logger.RTrace("Marshal json form of provider type " + serviceType)
	formString, err := json.Marshal(formData)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	logger.RTrace("Marshal json form done")
	ctx.JSON(http.StatusOK, dtos.ServiceStoreItemFormDto{FormJson: string(formString)})
}
