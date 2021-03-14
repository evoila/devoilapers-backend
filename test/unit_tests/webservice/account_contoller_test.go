package webservice

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	provider2 "OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_AccountController_HandlePostLogin(t *testing.T) {
	var provider provider2.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "TestType"
		},
	}

	router := CreateRouter(t, &provider)

	// Valid credentials
	requestDto := dtos.AccountCredentialsDto{
		Password: TEST_PASSWORD,
		Username: TEST_USERNAME,
	}

	var responseDto dtos.AuthenticationResponseDataDto
	statusCode := MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/accounts/login",
		&requestDto,
		&responseDto,
	)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, responseDto.IsValid)
	assert.Equal(t, TEST_ROLE, responseDto.Role)

	// Invalid credentials
	requestDto = dtos.AccountCredentialsDto{
		Password: TEST_PASSWORD + "X",
		Username: TEST_USERNAME,
	}

	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/accounts/login",
		&requestDto,
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// Invalid body
	errDto := dtos.HTTPErrorDto{}
	statusCode = MakeRequest(
		t,
		router,
		false,
		http.MethodPost,
		"/api/v1/accounts/login",
		nil,
		&errDto,
	)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, http.StatusBadRequest, errDto.Code)
	assert.NotEqual(t, "", errDto.Message)
}
