package controller

import (
	"OperatorAutomation/cmd/service/dtos"
	"OperatorAutomation/cmd/service/utils"
	"github.com/gin-gonic/gin"
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
			Type:        (*provider).GetName(),
			Description: (*provider).GetServiceDescription(),
			ImageBase64: (*provider).GetServiceImage(),
		}

		serviceStoreOverviewData.ServiceStoreItems = append(serviceStoreOverviewData.ServiceStoreItems, serviceStoreItem)
	}

	ctx.JSON(http.StatusOK, serviceStoreOverviewData)
}

// Default Service Yaml-Template godoc
// @Summary Get the default yaml for a service-template
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
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /servicestore/yaml/{servicetype} [get]
func (controller ServiceStoreController) HandleGetServiceStoreItemYaml(ctx *gin.Context) {
	serviceType := ctx.Param("servicetype")

	provider, err := controller.Core.GetProviderByName(serviceType)
	if err != nil {
		utils.NewError(ctx, http.StatusNotFound, err)
		return
	}

	serviceYaml := dtos.ServiceStoreItemYamlDto{
		Yaml: (*(*provider).GetTemplate()).GetYAML(),
	}

	ctx.JSON(http.StatusOK, serviceYaml)
}
