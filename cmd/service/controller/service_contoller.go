package controller

import (
	"OperatorAutomation/cmd/service/dtos"
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/pkg/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
func HandlePostCreateServiceInstance(ctx *gin.Context, core *core.Core) {
	var yamlData dtos.ServiceYamlDto
	if err := ctx.ShouldBindJSON(&yamlData); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update service instance godoc
// @Summary Update service instance from yaml
// @Description Update an instance of a service from yaml
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param serviceyaml body dtos.ServiceYamlDto true "Service-Yaml"
// @Param serviceid path string true "Id of service"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/update/{serviceid} [post]
func HandlePostUpdateServiceInstance(ctx *gin.Context, c *core.Core) {
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
// @Param serviceid path string true "Id of service"
// @Param actioncommand path string true "action command"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/action/{serviceid}/{actioncommand}  [post]
func HandlePostServiceInstanceAction(ctx *gin.Context, c *core.Core) {

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
// @Param serviceid path string true "Id of service"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/{serviceid} [delete]
func HandleDeleteServiceInstance(ctx *gin.Context, c *core.Core) {
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
// @Param serviceid path string true "Id of service"
//
// @Success 200 {object} dtos.ServiceInstanceDetailsOverviewDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/info/{serviceid} [get]
func HandleGetServiceInstanceDetails(ctx *gin.Context, c *core.Core) {
	serviceId := ctx.Param("serviceid")

	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{
			dtos.ServiceInstanceDetailsDto{
				Name:      "Instance 1",
				Id:        serviceId,
				Type:      "kibana",
				Status:    "ok",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					dtos.ServiceInstanceActionGroupDto{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							dtos.ServiceInstanceActionDto{
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
func HandleGetServiceInstanceDetailsForAllInstances(ctx *gin.Context, c *core.Core) {

	// Return all instances
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{
			dtos.ServiceInstanceDetailsDto{
				Name:      "Instance 1",
				Id:        "111111111-1111-4D9D-80C7-02AF85C822A8",
				Type:      "kibana",
				Status:    "ok",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					dtos.ServiceInstanceActionGroupDto{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							dtos.ServiceInstanceActionDto{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
			dtos.ServiceInstanceDetailsDto{
				Name:      "Instance 2",
				Id:        "22222222-XXXX-4DDD-80C7-02AF85999999",
				Type:      "elasticsearch",
				Status:    "warning",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					dtos.ServiceInstanceActionGroupDto{
						GroupName: "Backup and Restore",
						Actions: []dtos.ServiceInstanceActionDto{
							dtos.ServiceInstanceActionDto{
								Name:    "Backup",
								Command: "cmd_backup_elasticsearch",
							},
							dtos.ServiceInstanceActionDto{
								Name:    "Restore",
								Command: "cmd_restore_elasticsearch",
							},
						},
					},
					dtos.ServiceInstanceActionGroupDto{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							dtos.ServiceInstanceActionDto{
								Name:    "Expose",
								Command: "cmd_expose",
							},
						},
					},
				},
			},
			dtos.ServiceInstanceDetailsDto{
				Name:      "Instance 3",
				Id:        "33333333-XXXX-4DDD-80C7-02AF85999999",
				Type:      "logstash",
				Status:    "error",
				Namespace: "user_namespace_42",
				ActionGroups: []dtos.ServiceInstanceActionGroupDto{
					dtos.ServiceInstanceActionGroupDto{
						GroupName: "Security",
						Actions: []dtos.ServiceInstanceActionDto{
							dtos.ServiceInstanceActionDto{
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
// @Description Get the yaml file for an specific service instance. Parameter serviceid has to be supplied.
//
// @tags Service
//
// @Accept json
// @Produce json
//
// @Security BasicAuth
//
// @Param serviceid path string true "Id of service"
//
// @Success 200 {object} dtos.ServiceYamlDto
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/info/{serviceid} [get]
func HandleGetServiceInstanceYaml(ctx *gin.Context, c *core.Core) {
	serviceId := ctx.Param("serviceid")

	yamlData := dtos.ServiceYamlDto{
		Yaml: "item: " + serviceId,
	}

	ctx.JSON(http.StatusOK, yamlData)
}
