package controller

import (
	"OperatorAutomation/cmd/service/dtos"
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/pkg/core/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	serviceType := ctx.Param("servicetype")
	userCtx := controller.Core.CrateUserContext(userInfos)
	err := userCtx.CreateServices(serviceType, yamlData.Yaml)

	if err != nil {
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
// @Param payload body string true "Payload"
// @Param servicetype path string true "Type of service"
// @Param servicename path string true "Id of service"
// @Param actioncommand path string true "action command"
//
// @Success 200
// @Failure 401 {object} dtos.HTTPErrorDto
//
// @Router /services/action/{servicetype}/{servicename}/{actioncommand}  [post]
func (controller ServiceController) HandlePostServiceInstanceAction(ctx *gin.Context) {

	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")
	serviceActionCommand := ctx.Param("actioncommand")

	userCtx := controller.Core.CrateUserContext(userInfos)
	service, err := userCtx.GetService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	jsonData, err := ioutil.ReadAll(ctx.Request.Body)

	for _, group := range (*service).GetActions() {
		for _, action := range group.GetActions() {
			if action.GetUniqueCommand() != serviceActionCommand {
				continue
			}

			placeholder := action.GetPlaceholder()
			if err := json.Unmarshal(jsonData, placeholder); err != nil {
				utils.NewError(ctx, http.StatusBadRequest, err)
				return
			}

			val, err := action.GetActionExecuteCallback()(placeholder)
			if err != nil {
				utils.NewError(ctx, http.StatusBadRequest, err)
				return
			}

			ctx.JSON(http.StatusOK, val)
			return
		}
	}

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
	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	userCtx := controller.Core.CrateUserContext(userInfos)
	err := userCtx.DeleteService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

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
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	userCtx := controller.Core.CrateUserContext(userInfos)
	servicePtr, err := userCtx.GetService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	service := *servicePtr
	serviceDto := dtos.ServiceInstanceDetailsDto{}
	serviceDto.Status = serviceStatusToString(service.GetStatus())
	serviceDto.Type = serviceType
	serviceDto.Name = service.GetName()
	serviceDto.ActionGroups = serviceGroupToDto(servicePtr)
	instanceDetailsOverview.Instances = append(instanceDetailsOverview.Instances, serviceDto)

	ctx.JSON(http.StatusOK, instanceDetailsOverview)
}

func serviceGroupToDto(servicePtr *service.IService) []dtos.ServiceInstanceActionGroupDto {
	actionGroups := []dtos.ServiceInstanceActionGroupDto{}
	service := *servicePtr

	for _, group := range service.GetActions() {
		groupDto := dtos.ServiceInstanceActionGroupDto{Actions: []dtos.ServiceInstanceActionDto{}}
		groupDto.GroupName = group.GetName()

		for _, action := range group.GetActions() {

			jsonPlaceholder, _ := json.Marshal(action.GetPlaceholder())

			actionDto := dtos.ServiceInstanceActionDto{
				Name: action.GetName(),
				Command: action.GetUniqueCommand(),
				Placeholder: string(jsonPlaceholder),
			}

			groupDto.Actions = append(groupDto.Actions, actionDto)
		}

		actionGroups = append(actionGroups, groupDto)
	}

	return actionGroups
}

func serviceStatusToString(status int) string {
	switch status {
	case 0:
		return "Error"
	case 1:
		return "Warning"
	case 2:
		return "Ok"
	default:
		return "Unkown"
	}
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
	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	userCtx := controller.Core.CrateUserContext(userInfos)
	services := userCtx.GetServices()

	for _, servicePtr := range services {
		service := *servicePtr
		serviceDto := dtos.ServiceInstanceDetailsDto{}
		serviceDto.Status = serviceStatusToString(service.GetStatus())
		serviceDto.Type = service.GetType()
		serviceDto.Name = service.GetName()
		serviceDto.ActionGroups = serviceGroupToDto(servicePtr)
		instanceDetailsOverview.Instances = append(instanceDetailsOverview.Instances, serviceDto)
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
	servicetype := ctx.Param("servicetype")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos, foundUser := controller.UserManagement.GetUserInformation(user, password)
	if !foundUser {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	userCtx := controller.Core.CrateUserContext(userInfos)
	service, err := userCtx.GetService(servicetype, servicename)
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	yamlData := dtos.ServiceYamlDto{
		Yaml: (*service).GetTemplate().GetYAML(),
	}

	ctx.JSON(http.StatusOK, yamlData)
}
