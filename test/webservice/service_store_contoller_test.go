package webservice

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/controller"
	"OperatorAutomation/cmd/service/dtos"
	user "OperatorAutomation/cmd/service/management"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/test/common_test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Create_Users(t *testing.T) []config.User {
	return []config.User{
		{
			Name:     "Test",
			Password: "Test",
		},
	}
}

func Create_Base_Controller_And_Gin_Ctx(provider *service.IServiceProvider, t *testing.T) (*gin.Context, controller.BaseController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	core1 := core.CreateCore([]*service.IServiceProvider{provider})
	userManagement := user.CreateUserManagement(Create_Users(t))

	baseController := controller.BaseController{Core: core1, UserManagement: userManagement}

	return c, baseController, w
}

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

	ginCtx, baseController, responseWriter := Create_Base_Controller_And_Gin_Ctx(&provider, t)
	serviceStoreController := controller.ServiceStoreController{BaseController: baseController}
	serviceStoreController.HandleGetServiceStoreOverview(ginCtx)

	var dto dtos.ServiceStoreOverviewDto
	common_test.GetMessageAndCode(responseWriter, &dto)

	assert.Equal(t, 1, len(dto.ServiceStoreItems))
	serviceStoreItem := dto.ServiceStoreItems[0]
	assert.Equal(t, "TestType", serviceStoreItem.Type)
	assert.Equal(t, "TestDescription", serviceStoreItem.Description)
	assert.Equal(t, "TestImage", serviceStoreItem.ImageBase64)
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
		GetDescriptionCb: func() string {
			return "TestDescription"
		},
		GetServiceImageCb: func() string {
			return "TestImage"
		},
	}

	ginCtx, baseController, responseWriter := Create_Base_Controller_And_Gin_Ctx(&provider, t)

	serviceStoreController := controller.ServiceStoreController{BaseController: baseController}

	router := gin.Default()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(responseWriter, req)

	serviceStoreController.HandleGetServiceStoreItemYaml(ginCtx)


	var dto dtos.ServiceStoreOverviewDto
	common_test.GetMessageAndCode(responseWriter, &dto)

	assert.Equal(t, 1, len(dto.ServiceStoreItems))
	serviceStoreItem := dto.ServiceStoreItems[0]
	assert.Equal(t, "TestType", serviceStoreItem.Type)
	assert.Equal(t, "TestDescription", serviceStoreItem.Description)
	assert.Equal(t, "TestImage", serviceStoreItem.ImageBase64)
}
