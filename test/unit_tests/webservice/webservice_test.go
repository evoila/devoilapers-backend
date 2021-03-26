package webservice

import (
	"OperatorAutomation/cmd/service/config"
	ws "OperatorAutomation/cmd/service/webserver"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/provider"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var TEST_USERNAME = "TestUsername"
var TEST_PASSWORD = "TestPassword"
var TEST_ROLE = "user"
var TEST_NAMESPACE = "TestNamespace"
var TEST_TOKEN = "TestToken"

func Test_Webservice(t *testing.T) {

}

func CreateRouter(t *testing.T, newProvider *provider.IServiceProvider) *gin.Engine {
	core1 := core.CreateCore([]*provider.IServiceProvider{newProvider})
	router := ws.BuildRouter(Create_AppConfig(t), core1)
	return router
}

func MakeRequest(t *testing.T,
	router *gin.Engine,
	doAuth bool,
	httpMethod string,
	path string,
	bodyObj interface{},
	responseObj interface{}) int {

	// Append a json body only if bodyObj is set
	var payload io.Reader = nil
	if bodyObj != nil {
		serializedTargetBytes, err := json.Marshal(bodyObj)
		if err != nil {
			t.FailNow()
		}

		payload = strings.NewReader(string(serializedTargetBytes))
	}

	request, _ := http.NewRequest(httpMethod, path, payload)

	// May append auth infos to be able to make the request
	if doAuth {
		request.SetBasicAuth(TEST_USERNAME, TEST_PASSWORD)
	}

	// Execute request
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseObj != nil {
		// Deserialize method response
		err := json.Unmarshal(responseRecorder.Body.Bytes(), responseObj)

		// Let the test fail if the response object is not what was expected
		if err != nil {
			t.FailNow()
		}
	} else {
		// Ensure that if not object is expected the server should not return anything
		assert.Equal(t, 0, len(responseRecorder.Body.Bytes()))
	}

	return responseRecorder.Code
}

func Create_Users(t *testing.T) []config.User {
	return []config.User{
		{
			Name:                  TEST_USERNAME,
			Password:              TEST_PASSWORD,
			KubernetesNamespace:   TEST_NAMESPACE,
			KubernetesAccessToken: TEST_TOKEN,
		},
	}
}

func Create_AppConfig(t *testing.T) config.RawConfig {
	return config.RawConfig{
		Users: Create_Users(t),
	}
}
