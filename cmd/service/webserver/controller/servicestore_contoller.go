package controller

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/cmd/service/utils"
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
	serviceStoreOverviewData := dtos.ServiceStoreOverviewDto{ServiceStoreItems: []dtos.ServiceStoreItemDto{}}

	for _, provider := range controller.Core.Providers {
		serviceStoreItem := dtos.ServiceStoreItemDto{
			Type:        (*provider).GetServiceType(),
			Description: (*provider).GetServiceDescription(),
			ImageBase64: (*provider).GetServiceImage(),
		}

		serviceStoreOverviewData.ServiceStoreItems = append(serviceStoreOverviewData.ServiceStoreItems, serviceStoreItem)
	}

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
	serviceType := ctx.Param("servicetype")

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

	yamlString, err := yaml.Marshal(yamlObject)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Build response
	serviceYaml := dtos.ServiceStoreItemYamlDto{
		Yaml: string(yamlString),
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
// @Success 200 {object} dtos.ServiceStoreItemYamlDto
// @Failure 400 {object} dtos.HTTPErrorDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /servicestore/form/{servicetype} [get]
func (controller ServiceStoreController) HandleGetServiceStoreItemForm(ctx *gin.Context) {
	serviceType := ctx.Param("servicetype")

	provider, err := controller.Core.GetProviderByName(serviceType)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	// Generate form data
	formData, err := (*provider).GetJsonForm(userInfos)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, formData)
}
