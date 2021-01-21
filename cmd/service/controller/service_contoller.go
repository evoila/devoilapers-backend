package controller

import (
	"OperatorAutomation/cmd/service/dtos"
	"OperatorAutomation/cmd/service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServiceController struct {
	BaseController
}

// Create service instance godoc
// @Summary Create service instance from yaml
// @Description Create an instance of a service from yaml
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param serviceyaml body dtos.ServiceYamlDto true "Service-Yaml"
// @Param servicetype path string true "Type of service"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/create/{servicetype} [post]
func (controller ServiceController) HandlePostCreateServiceInstance(ctx *gin.Context) {
	var yamlData dtos.ServiceYamlDto
	if err := ctx.ShouldBindJSON(&yamlData); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// Apply a service specific action godoc
// @Summary Apply a service specific action to a service instance
// @Description Apply a service specific action to a service instance
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param servicetype path string true "Type of service"
// @Param servicename path string true "Id of service"
// @Param actioncommand path string true "action command"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/action/{servicetype}/{servicename}/{actioncommand}  [post]
func (controller ServiceController) HandlePostServiceInstanceAction(ctx *gin.Context) {

	ctx.Status(http.StatusOK)
}

// Delete service instance godoc
// @Summary Delete a service instance
// @Description Delete an instance of a service
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param servicetype path string true "Type of service"
// @Param servicename path string true "Id of service"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/{servicetype}/{servicename} [delete]
func (controller ServiceController) HandleDeleteServiceInstance(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// Get service instance details godoc
// @Summary Get details over a single service instance
// @Description Get details over a single service instance
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param servicetype path string true "Type of service"
// @Param servicename path string true "Id of service"
//
// @Success 200 {object} dtos.ServiceInstanceDetailsOverviewDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/info/{servicetype}/{servicename} [get]
func (controller ServiceController) HandleGetServiceInstanceDetails(ctx *gin.Context) {
	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{
			{
				Name:      "Instance 1",
				Type:      "kibana",
				Status:    "ok",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
		},
	}

	ctx.JSON(http.StatusOK, instanceDetailsOverview)
	return

}

// Get service instance overview godoc
// @Summary Get an overview over all service instances
// @Description Get an overview over all service instances
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Success 200 {object} dtos.ServiceInstanceDetailsOverviewDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/info [get]
func (controller ServiceController) HandleGetServiceInstanceDetailsForAllInstances(ctx *gin.Context) {

	// Return all instances
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{
			{
				Name:      "Instance 1",
				Type:      "kibana",
				Status:    "ok",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
			{
				Name:      "Instance 2",
				Type:      "elasticsearch",
				Status:    "warning",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					{
						GroupName: "Backup and Restore",
						Actions: []dtos.ServiceInstanceActionDto{
							{
								Name:    "Backup",
								Command: "cmd_backup_elasticsearch",
							},
							{
								Name:    "Restore",
								Command: "cmd_restore_elasticsearch",
							},
						},
					},
					{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
			{
				Name:      "Instance 3",
				Type:      "logstash",
				Status:    "error",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
		},
	}

	ctx.JSON(http.StatusOK, instanceDetailsOverview)
	return
}

// Get service instance yaml godoc
// @Summary Get the yaml file for an instance
// @Description Get the yaml file for an specific service instance. Parameter servicename has to be supplied.
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param servicetype path string true "Type of service"
// @Param servicename path string true "Id of service"
//
// @Success 200 {object} dtos.ServiceYamlDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/yaml/{servicetype}/{servicename} [get]
func (controller ServiceController) HandleGetServiceInstanceYaml(ctx *gin.Context) {
	servicename := ctx.Param("servicename")

	yamlData := dtos.ServiceYamlDto{
		Yaml: "item: " + servicename,
	}

	ctx.JSON(http.StatusOK, yamlData)
}
