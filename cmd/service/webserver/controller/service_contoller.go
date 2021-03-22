package controller

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/pkg/core/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ServiceController struct {
	BaseController
}

const ServiceLogPrefix = "File: service_controller.go: "

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
	log.Trace(ServiceLogPrefix + "Received post request to create an instance of a service.")

	log.Trace(ServiceLogPrefix + "Going to unmarshal body.")
	var yamlData dtos.ServiceYamlDto
	if err := ctx.ShouldBindJSON(&yamlData); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	log.Trace(ServiceLogPrefix + "Going to create a usercontext.")
	userCtx := controller.Core.CrateUserContext(userInfos)

	log.Trace(ServiceLogPrefix + "Going to create a service of type " + serviceType + " out of the given yaml.")
	err := userCtx.CreateServices(serviceType, yamlData.Yaml)

	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	log.Trace(ServiceLogPrefix + "Creation of service done.")
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
	log.Trace(ServiceLogPrefix + "Received post request to execute an action.")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")
	serviceActionCommand := ctx.Param("actioncommand")

	log.Trace(ServiceLogPrefix + "Action command " + serviceActionCommand +
		" should be executed on service " + serviceType + "/" + serviceName + ".")

	log.Trace(ServiceLogPrefix + "Going to create a user context.")
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

	log.Trace(ServiceLogPrefix + "Going to parse body.")
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)

	for _, group := range (*service).GetActionGroups() {
		for _, action := range group.GetActions() {
			if action.GetUniqueCommand() != serviceActionCommand {
				continue
			}

			log.Trace(ServiceLogPrefix + "Action " + serviceActionCommand + " could be found.")
			log.Trace(ServiceLogPrefix + "Getting placeholder of action.")

			placeholder := action.GetJsonFormResultPlaceholder()
			if placeholder != nil {
				log.Trace(ServiceLogPrefix + "Filling placeholder of action with given data.")
				if err := json.Unmarshal(jsonData, placeholder); err != nil {
					utils.NewError(ctx, http.StatusBadRequest, err)
					return
				}
			}

			log.Trace(ServiceLogPrefix + "Executing the action.")
			actionResult, err := action.GetActionExecuteCallback()(placeholder)
			if err != nil {
				utils.NewError(ctx, http.StatusInternalServerError, err)
				return
			}


			log.Trace(ServiceLogPrefix + "Parsing the result of action.")
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

	log.Warn(ServiceLogPrefix + "Action " + serviceActionCommand + " could not be found.")
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
	log.Trace(ServiceLogPrefix + "Received request to delete a service instance.")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	log.Trace(ServiceLogPrefix + "Going to create a user context.")
	userCtx := controller.Core.CrateUserContext(userInfos)

	log.Trace(ServiceLogPrefix + "Going to delete service of type " + serviceType + " with name " + serviceName + ".")
	err := userCtx.DeleteService(serviceType, serviceName)

	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	log.Trace(ServiceLogPrefix + "Service deleted.")
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
	log.Trace(ServiceLogPrefix + "Received request to get service instance details.")

	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	serviceType := ctx.Param("servicetype")
	serviceName := ctx.Param("servicename")

	log.Trace(ServiceLogPrefix + "Going to create a user context.")
	userCtx := controller.Core.CrateUserContext(userInfos)
	log.Trace(ServiceLogPrefix + "Get service of type " + serviceType + " with name " + serviceName + ".")
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

	log.Trace(ServiceLogPrefix + "Service could be found")
	ctx.JSON(http.StatusOK, instanceDetailsOverview)
}

func serviceActionGroupToDto(servicePtr *service.IService) []dtos.ServiceInstanceActionGroupDto {
	log.Trace(ServiceLogPrefix + "Converting service action groups to dtos.")
	actionGroups := []dtos.ServiceInstanceActionGroupDto{}
	service := *servicePtr

	for _, group := range service.GetActionGroups() {
		log.Trace(ServiceLogPrefix + "Found group with name " + group.GetName() + ".")

		groupDto := dtos.ServiceInstanceActionGroupDto{Actions: []dtos.ServiceInstanceActionDto{}}
		groupDto.GroupName = group.GetName()

		for _, action := range group.GetActions() {
			log.Trace(ServiceLogPrefix + "Found action with name " + action.GetName() +
				" in group with name " + group.GetName() + ".")

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
	log.Trace(ServiceLogPrefix + "Parsing numeric error to string.")

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
	log.Trace(ServiceLogPrefix + "Received get request to get details for all service instances.")

	//Return single instance
	instanceDetailsOverview := dtos.ServiceInstanceDetailsOverviewDto{
		Instances: []dtos.ServiceInstanceDetailsDto{},
	}

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	log.Trace(ServiceLogPrefix + "Going to create a user context.")
	userCtx := controller.Core.CrateUserContext(userInfos)
	services, err := userCtx.GetServices()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for _, servicePtr := range services {
		service := *servicePtr
		log.Trace(ServiceLogPrefix + "Found service with name " + service.GetName() +
			" of type " + service.GetType() + ".")

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
	log.Trace(ServiceLogPrefix + "Received get request to get yaml of a service")

	servicename := ctx.Param("servicename")
	servicetype := ctx.Param("servicetype")

	user, password, _ := ctx.Request.BasicAuth()
	userInfos := controller.UserManagement.GetUserInformation(user, password)

	log.Trace(ServiceLogPrefix + "Going to create a user context.")
	userCtx := controller.Core.CrateUserContext(userInfos)

	log.Trace(ServiceLogPrefix + "Get service with name " + servicename + " of type " + servicetype)
	service, err := userCtx.GetService(servicetype, servicename)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	log.Trace(ServiceLogPrefix + "Service found. Getting yaml from it.")
	yamlData := dtos.ServiceYamlDto{
		Yaml: (*service).GetYamlTemplate(),
	}

	ctx.JSON(http.StatusOK, yamlData)
}
