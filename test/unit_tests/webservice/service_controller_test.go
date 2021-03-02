package webservice

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/test/unit_tests/common_test"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_ServiceController_HandlePostCreateServiceInstance(t *testing.T) {

	createServiceGotCalled := 0
	providerError := false
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		CreateServiceCb: func(auth common.IKubernetesAuthInformation, yaml string) error {
			assert.Equal(t, TEST_TOKEN, auth.GetKubernetesAccessToken())
			assert.Equal(t, TEST_NAMESPACE, auth.GetKubernetesNamespace())
			assert.Equal(t, "SomeYaml", yaml)

			createServiceGotCalled += 1

			if providerError {
				return errors.New("SomeError")
			}

			return nil
		},
	}

	router := CreateRouter(t, &provider)

	// Valid credentials
	requestDto := dtos.ServiceYamlDto{
		Yaml: "SomeYaml",
	}

	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/create/"+provider.GetServiceType(),
		&requestDto,
		nil,
	)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, 1, createServiceGotCalled)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/services/create/"+provider.GetServiceType(),
		&requestDto,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, 1, createServiceGotCalled)

	// On request error
	var errorDto dtos.HTTPErrorDto
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/create/"+provider.GetServiceType(),
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, 1, createServiceGotCalled)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)

	// On provider error
	providerError = true
	errorDto = dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/create/"+provider.GetServiceType(),
		&requestDto,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, 2, createServiceGotCalled)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceController_HandlePostServiceInstanceAction(t *testing.T) {

	actionGotCalled := 0
	actionError := false
	serviceError := false

	// Create provider
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetServiceCb: func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
			var service0 service.IService = common_test.TestService{
				GetNameCb: func() string {
					return "TestService"
				},
				GetActionGroupCb: func() []action.IActionGroup {
					return []action.IActionGroup{
						action.ActionGroup{
							Name: "TestGroup",
							Actions: []action.IAction{
								action.Action{
									Name:          "TestAction",
									Placeholder:   &common_test.TestPlaceholder{},
									UniqueCommand: "TestActionCmd",
									ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
										actionGotCalled += 1
										if actionError {
											return placeholder, errors.New("Some error")
										}

										return placeholder, nil
									},
								},
							},
						},
					}
				},
			}

			if serviceError {
				return nil, errors.New("Some error")
			}

			return &service0, nil
		},
	}

	router := CreateRouter(t, &provider)

	// Valid credentials
	requestDto := common_test.TestPlaceholder{SomeValue: "MyTestValue"}
	responseDto := common_test.TestPlaceholder{}

	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/action/TestType/TestService/TestActionCmd",
		&requestDto,
		&responseDto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, requestDto.SomeValue, responseDto.SomeValue)
	assert.Equal(t, 1, actionGotCalled)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/services/action/TestType/TestService/TestActionCmd",
		&requestDto,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, 1, actionGotCalled)

	// On request error
	var errorDto dtos.HTTPErrorDto
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/action/TestType/TestService/TestActionCmd",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, 1, actionGotCalled)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)

	// On action error
	actionError = true
	errorDto = dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/action/TestType/TestService/TestActionCmd",
		&requestDto,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, 2, actionGotCalled)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)

	// On service error
	serviceError = true
	errorDto = dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/services/action/TestType/TestService/TestActionCmd",
		&requestDto,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, 2, actionGotCalled)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceController_HandleDeleteServiceInstance(t *testing.T) {

	deleteGotCalled := 0
	deleteError := false
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		DeleteServiceCb: func(auth common.IKubernetesAuthInformation, id string) error {
			assert.Equal(t, "TestService", id)

			deleteGotCalled += 1

			if deleteError {
				return errors.New("SomeError")
			}

			return nil
		},
	}

	router := CreateRouter(t, &provider)

	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodDelete,
		"/api/v1/services/TestType/TestService",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 1, deleteGotCalled)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodDelete,
		"/api/v1/services/TestType/TestService",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, 1, deleteGotCalled)

	// On provider error
	deleteError = true
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodDelete,
		"/api/v1/services/TestType/TestService",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, 2, deleteGotCalled)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceController_HandleGetServiceInstanceDetails(t *testing.T) {
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetServiceCb: func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
			assert.Equal(t, "TestService", id)

			var service0 service.IService = common_test.TestService{
				GetStatusCb: func() int {
					return service.SERVICE_STATUS_WARNING
				},
				GetNameCb: func() string {
					return id
				},
				GetTypeCb: func() string {
					return "TestType"
				},
				GetActionGroupCb: func() []action.IActionGroup {
					return []action.IActionGroup{
						action.ActionGroup{
							Name: "SomeActionGroup",
							Actions: []action.IAction{
								action.Action{
									Name: "MyAction",
									Placeholder: &common_test.TestPlaceholder{},
									UniqueCommand: "MyActionCmd",
								},
							},
						},
					}
				},
			}

			return &service0, nil
		},
	}

	router := CreateRouter(t, &provider)

	// Ok
	responseDto := dtos.ServiceInstanceDetailsOverviewDto{}
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/info/TestType/TestService",
		nil,
		&responseDto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 1, len(responseDto.Instances))
	responseInstance := responseDto.Instances[0]
	assert.Equal(t, "TestService", responseInstance.Name)
	assert.Equal(t, "TestType", responseInstance.Type)
	assert.Equal(t, "Warning", responseInstance.Status)
	assert.Equal(t, 1, len(responseInstance.ActionGroups))
	responseActionGroup := responseInstance.ActionGroups[0]
	assert.Equal(t, "SomeActionGroup", responseActionGroup.GroupName)
	assert.Equal(t, 1, len(responseActionGroup.Actions))
	responseAction := responseActionGroup.Actions[0]
	assert.Equal(t, "MyAction", responseAction.Name)
	assert.Equal(t, "MyActionCmd", responseAction.Command)
	assert.Equal(t, "{\"SomeValue\":\"\"}", responseAction.Placeholder)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/services/info/TestType/TestService",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// On provider error
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/info/TestTypeNotExising/TestService",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceController_HandleGetServiceInstanceDetailsForAllInstances(t *testing.T) {

	providerError := false
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetServicesCb: func(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
			var service0 service.IService = common_test.TestService{
				GetStatusCb: func() int {
					return service.SERVICE_STATUS_ERROR
				},
				GetNameCb: func() string {
					return "TestService0"
				},
				GetTypeCb: func() string {
					return "TestType"
				},
				GetActionGroupCb: func() []action.IActionGroup {
					return []action.IActionGroup{
						action.ActionGroup{
							Name: "SomeActionGroup0",
							Actions: []action.IAction{
								action.Action{
									Name: "MyAction0",
									Placeholder: &common_test.TestPlaceholder{},
									UniqueCommand: "MyActionCmd0",
								},
							},
						},
					}
				},
			}

			var service1 service.IService = common_test.TestService{
				GetStatusCb: func() int {
					return service.SERVICE_STATUS_OK
				},
				GetNameCb: func() string {
					return "TestService1"
				},
				GetTypeCb: func() string {
					return "TestType"
				},
				GetActionGroupCb: func() []action.IActionGroup {
					return []action.IActionGroup{
						action.ActionGroup{
							Name: "SomeActionGroup1",
							Actions: []action.IAction{
								action.Action{
									Name: "MyAction1",
									Placeholder: &common_test.TestPlaceholder{},
									UniqueCommand: "MyActionCmd1",
								},
							},
						},
					}
				},
			}

			if providerError == true {
				return nil, errors.New("Some error")
			}

			return []*service.IService{&service0, &service1}, nil
		},
	}


	router := CreateRouter(t, &provider)

	// Ok
	responseDto := dtos.ServiceInstanceDetailsOverviewDto{}
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/info",
		nil,
		&responseDto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 2, len(responseDto.Instances))

	responseInstance0 := responseDto.Instances[0]
	assert.Equal(t, "TestService0", responseInstance0.Name)
	assert.Equal(t, "TestType", responseInstance0.Type)
	assert.Equal(t, "Error", responseInstance0.Status)
	assert.Equal(t, 1, len(responseInstance0.ActionGroups))

	responseInstance1 := responseDto.Instances[1]
	assert.Equal(t, "TestService1", responseInstance1.Name)
	assert.Equal(t, "TestType", responseInstance1.Type)
	assert.Equal(t, "Ok", responseInstance1.Status)
	assert.Equal(t, 1, len(responseInstance1.ActionGroups))

	responseActionGroup0 := responseInstance0.ActionGroups[0]
	assert.Equal(t, "SomeActionGroup0", responseActionGroup0.GroupName)
	assert.Equal(t, 1, len(responseActionGroup0.Actions))
	responseAction0 := responseActionGroup0.Actions[0]

	responseActionGroup1 := responseInstance1.ActionGroups[0]
	assert.Equal(t, "SomeActionGroup1", responseActionGroup1.GroupName)
	assert.Equal(t, 1, len(responseActionGroup1.Actions))
	responseAction1 := responseActionGroup1.Actions[0]

	assert.Equal(t, "MyAction0", responseAction0.Name)
	assert.Equal(t, "MyActionCmd0", responseAction0.Command)
	assert.Equal(t, "{\"SomeValue\":\"\"}", responseAction0.Placeholder)

	assert.Equal(t, "MyAction1", responseAction1.Name)
	assert.Equal(t, "MyActionCmd1", responseAction1.Command)
	assert.Equal(t, "{\"SomeValue\":\"\"}", responseAction1.Placeholder)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/services/info",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// On provider error
	errorDto := dtos.HTTPErrorDto{}
	providerError = true
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/info",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceController_HandleGetServiceInstanceYaml(t *testing.T) {
	providerError := false
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetServiceCb: func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
			assert.Equal(t, "TestService", id)

			var service0 service.IService = common_test.TestService{
				GetTemplateCb: func() service.IServiceTemplate {
					return service.ServiceTemplate{
						Yaml: "TestYaml",
					}
				},
			}

			if providerError {
				return nil, errors.New("Some error")
			}

			return &service0, nil
		},
	}

	router := CreateRouter(t, &provider)

	// Ok
	responseDto := dtos.ServiceYamlDto{}
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/yaml/TestType/TestService",
		nil,
		&responseDto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "TestYaml", responseDto.Yaml)

	// Invalid credentials
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/services/yaml/TestType/TestService",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// On provider error
	providerError = true
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/services/yaml/TestType/TestService",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, http.StatusInternalServerError, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}