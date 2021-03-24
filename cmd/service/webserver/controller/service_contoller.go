package controller

import (
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/utils/logger"
	"encoding/json"
	"errors"
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
// @Success 201
// @Failure 400 {object} dtos.HTTPErrorDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/create/{servicetype} [post]
func (controller ServiceController) HandlePostCreateServiceInstance(ctx *gin.Context) {
	logger.RTrace("Received post request to create an instance of a service")
	
	logger.RTrace("Going to unmarshal body")
	var yamlData dtos.ServiceYamlDto
	if err := ctx.ShouldBindJSON(&yamlData); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	logger.RTrace("Going to create a usercontext")
	userCtx := controller.Core.CrateUserContext(userInfos)

	logger.RTrace("Going to create a service of type " + serviceType + " out of the given yaml")
	err := userCtx.CreateServices(serviceType, yamlData.Yaml)

	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	logger.RTrace("Creation of service done")
	ctx.Status(http.StatusCreated)
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
// @Success 200 {object} dtos.ServiceInstanceActionResponseDto
// @Failure 400 {object} dtos.HTTPErrorDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/action/{servicetype}/{servicename}/{actioncommand}  [post]
func (controller ServiceController) HandlePostServiceInstanceAction(ctx *gin.Context) {
	logger.RTrace("Received post request to execute an action")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")
	serviceActionCommand := ctx.Param("actioncommand")

	logger.RTrace("Action command " + serviceActionCommand +
		" should be executed on service " + serviceType + "/" + serviceName)

	logger.RTrace("Going to create a user context")
	userCtx := controller.Core.CrateUserContext(userInfos)
	service, err := userCtx.GetService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if ctx.Request.ContentLength == 0 {
		utils.NewError(ctx, http.StatusBadRequest, errors.New("Request does not contain any body data"))
		return
	}

	logger.RTrace("Going to parse body")
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)

	for _, group := range (*service).GetActionGroups() {
		for _, action := range group.GetActions() {
			if action.GetUniqueCommand() != serviceActionCommand {
				continue
			}

			logger.RTrace("Action " + serviceActionCommand + " could be found")
			logger.RTrace("Getting placeholder of action")

			placeholder := action.GetJsonFormResultPlaceholder()
			if placeholder != nil {
				logger.RTrace("Filling placeholder of action with given data")
				if err := json.Unmarshal(jsonData, placeholder); err != nil {
					utils.NewError(ctx, http.StatusBadRequest, err)
					return
				}
			}

			logger.RTrace("Executing the action")
			actionResult, err := action.GetActionExecuteCallback()(placeholder)
			if err != nil {
				utils.NewError(ctx, http.StatusInternalServerError, err)
				return
			}


			logger.RTrace("Parsing the result of action")
			actionResultJson := ""
			if actionResult != nil {
				actionResultBytes, err := json.Marshal(actionResult)
				if err != nil {
					utils.NewError(ctx, http.StatusInternalServerError, err)
					return
				}

				actionResultJson = string(actionResultBytes)
			}

			ctx.JSON(http.StatusOK, dtos.ServiceInstanceActionResponseDto{ResultJson: actionResultJson})
			return
		}
	}

	logger.RWarn("Action " + serviceActionCommand + " could not be found")
	utils.NewError(ctx, http.StatusBadRequest, errors.New("action " + serviceActionCommand + " could not be found"))
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
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/{servicetype}/{servicename} [delete]
func (controller ServiceController) HandleDeleteServiceInstance(ctx *gin.Context) {
	logger.RTrace("Received request to delete a service instance")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	logger.RTrace("Going to create a user context")
	userCtx := controller.Core.CrateUserContext(userInfos)

	logger.RTrace("Going to delete service of type " + serviceType + " with name " + serviceName)
	err := userCtx.DeleteService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	logger.RTrace("Service deleted")
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
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/info/{servicetype}/{servicename} [get]
func (controller ServiceController) HandleGetServiceInstanceDetails(ctx *gin.Context) {
	logger.RTrace("Received request to get service instance details")

	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	logger.RTrace("Going to create a user context")
	userCtx := controller.Core.CrateUserContext(userInfos)
	logger.RTrace("Get service of type " + serviceType + " with name " + serviceName)
	servicePtr, err := userCtx.GetService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	service := *servicePtr
	serviceDto := dtos.ServiceInstanceDetailsDto{}
	serviceDto.Status = serviceStatusToString(service.GetStatus())
	serviceDto.Type = serviceType
	serviceDto.Name = service.GetName()
	serviceDto.ActionGroups = serviceActionGroupToDto(servicePtr)
	instanceDetailsOverview.Instances = append(instanceDetailsOverview.Instances, serviceDto)

	logger.RTrace("Service could be found")
	ctx.JSON(http.StatusOK, instanceDetailsOverview)
}

func serviceActionGroupToDto(servicePtr *service.IService) []dtos.ServiceInstanceActionGroupDto {
	logger.RTrace("Converting service action groups to dtos")
	actionGroups := []dtos.ServiceInstanceActionGroupDto{}
	service := *servicePtr

	for _, group := range service.GetActionGroups() {
		logger.RTrace("Found group with name " + group.GetName())

		groupDto := dtos.ServiceInstanceActionGroupDto{Actions: []dtos.ServiceInstanceActionDto{}}
		groupDto.GroupName = group.GetName()

		for _, action := range group.GetActions() {
			logger.RTrace("Found action with name " + action.GetName() +
				" in group with name " + group.GetName() + " for " + service.GetType() +
				" service with name " + service.GetName())

			jsonPlaceholder, _ := json.Marshal(action.GetJsonForm())

			actionDto := dtos.ServiceInstanceActionDto{
				Name:     action.GetName(),
				Command:  action.GetUniqueCommand(),
				FormJson: string(jsonPlaceholder),
				IsToggle: action.GetIsToggleAction(),
			}

			groupDto.Actions = append(groupDto.Actions, actionDto)
		}

		actionGroups = append(actionGroups, groupDto)
	}

	return actionGroups
}

func serviceStatusToString(status int) string {
	logger.RTrace("Parsing numeric status code to string status")

	switch status {
	case service.ServiceStatusError:
		return "Error"
	case service.ServiceStatusWarning:
		return "Warning"
	case service.ServiceStatusOk:
		return "Ok"
	case service.ServiceStatusPending:
		return "Pending"
	default:
		return "Unknown"
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
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/info [get]
func (controller ServiceController) HandleGetServiceInstanceDetailsForAllInstances(ctx *gin.Context) {
	logger.RTrace("Received get request to get details for all service instances")

	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	logger.RTrace("Going to create a user context")
	userCtx := controller.Core.CrateUserContext(userInfos)
	services, err := userCtx.GetServices()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for _, servicePtr := range services {
		service := *servicePtr
		logger.RTrace("Found service with name " + service.GetName() +
			" of type " + service.GetType())

		serviceDto := dtos.ServiceInstanceDetailsDto{}
		serviceDto.Status = serviceStatusToString(service.GetStatus())
		serviceDto.Type = service.GetType()
		serviceDto.Name = service.GetName()
		serviceDto.ActionGroups = serviceActionGroupToDto(servicePtr)
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
// @Failure 500 {object} dtos.HTTPErrorDto
//
// @Router /services/yaml/{servicetype}/{servicename} [get]
func (controller ServiceController) HandleGetServiceInstanceYaml(ctx *gin.Context) {
	logger.RTrace("Received get request to get yaml of a service")

	servicename := ctx.Param("servicename")
	servicetype := ctx.Param("servicetype")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	logger.RTrace("Going to create a user context")
	userCtx := controller.Core.CrateUserContext(userInfos)

	logger.RTrace("Get service with name " + servicename + " of type " + servicetype)
	service, err := userCtx.GetService(servicetype, servicename)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	logger.RTrace("Service found. Getting yaml from it")
	yamlData := dtos.ServiceYamlDto{
		Yaml: (*service).GetYamlTemplate(),
	}

	ctx.JSON(http.StatusOK, yamlData)
}
