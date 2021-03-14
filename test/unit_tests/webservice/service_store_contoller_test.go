package webservice

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/pkg/core/common"
	provider2 "OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_ServiceStoreController_HandleGetServiceStoreOverview(t *testing.T) {
	var provider provider2.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetDescriptionCb: func() string {
			return "TestDescription"
		},
		GetServiceImageCb: func() string {
			return "TestImage"
		},
	}

	router := CreateRouter(t, &provider)

	// Authorized
	var dto dtos.ServiceStoreOverviewDto
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/servicestore/info",
		nil,
		&dto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 1, len(dto.ServiceStoreItems))
	serviceStoreItem := dto.ServiceStoreItems[0]
	assert.Equal(t, "TestType", serviceStoreItem.Type)
	assert.Equal(t, "TestDescription", serviceStoreItem.Description)
	assert.Equal(t, "TestImage", serviceStoreItem.ImageBase64)

	// Unauthorized
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/servicestore/info",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
}

func Test_ServiceStoreController_HandleGetServiceStoreForm(t *testing.T) {
	formTemplateError := false
	var provider provider2.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetFormTemplateCb: func(auth common.IKubernetesAuthInformation) (interface{}, error) {
			if formTemplateError {
				return nil, errors.New("My error")
			}

			return common_test.TestSerializableStruct{Value: "MyFormValue"}, nil
		},
	}

	router := CreateRouter(t, &provider)

	// Authorized
	var dto common_test.TestSerializableStruct
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/servicestore/form/"+provider.GetServiceType(),
		nil,
		&dto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "MyFormValue", dto.Value)

	// Unauthorized
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/servicestore/form/"+provider.GetServiceType(),
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// From template generation error
	formTemplateError = true
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/servicestore/form/"+provider.GetServiceType(),
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.NotEqual(t, "", errorDto.Message)
}

func Test_ServiceStoreController_HandleGetServiceStoreItemYaml(t *testing.T) {
	getYamlTemplateError := false
	var provider provider2.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetYamlTemplateCb: func(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
			filledForm := common_test.TestSerializableStruct{}
			err := json.Unmarshal(jsonFormResult, &filledForm)
			assert.Nil(t, err)

			if getYamlTemplateError {
				return nil, errors.New("Some error")
			}

			return common_test.TestSerializableStruct{Value: filledForm.Value}, nil
		},
	}

	router := CreateRouter(t, &provider)

	// Authorized
	filledForm := common_test.TestSerializableStruct{Value: "MyFilledFormValue"}
	var dto dtos.ServiceStoreItemYamlDto
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/servicestore/yaml/"+provider.GetServiceType(),
		filledForm,
		&dto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "value: MyFilledFormValue\n", dto.Yaml)

	// Not authorized
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/servicestore/yaml/"+provider.GetServiceType(),
		filledForm,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// Invalid provider
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/servicestore/yaml/NotExistingProvider",
		filledForm,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)

	// Invalid payload
	errorDto = dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/servicestore/yaml/NotExistingProvider",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)

	// Error during yaml generation
	getYamlTemplateError = true
	errorDto = dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodPost,
		"/api/v1/servicestore/yaml/"+provider.GetServiceType(),
		filledForm,
		&errorDto,
	)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.NotEqual(t, "", errorDto.Message)
}
