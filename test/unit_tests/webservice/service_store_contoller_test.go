package webservice

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_ServiceStoreController_HandleGetServiceStoreOverview(t *testing.T) {
	var provider service.IServiceProvider = common_test.TestProvider{
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

func Test_ServiceStoreController_HandleGetServiceStoreItemYaml(t *testing.T) {
	var provider service.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
		GetTemplateCb: func() *service.IServiceTemplate {
			var template service.IServiceTemplate = service.ServiceTemplate{
				Yaml: "TestYaml",
			}

			return &template
		},
	}

	router := CreateRouter(t, &provider)

	// Authorized
	var dto dtos.ServiceStoreItemYamlDto
	statusCode := MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/servicestore/yaml/" + provider.GetServiceType(),
		nil,
		&dto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "TestYaml", dto.Yaml)

	// Not authorized
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodGet,
		"/api/v1/servicestore/yaml/" + provider.GetServiceType(),
		nil,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)


	// Invalid provider
	errorDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		true,
		http.MethodGet,
		"/api/v1/servicestore/yaml/NotExistingProvider",
		nil,
		&errorDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, http.StatusBadRequest, errorDto.Code)
	assert.NotEqual(t, "", errorDto.Message)
}
