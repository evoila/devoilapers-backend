package controller

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"OperatorAutomation/cmd/service/utils"
	"OperatorAutomation/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountController struct {
	BaseController
}

// Login godoc
// @Summary User login
// @Description Get login token and role by providing username and password
//
// @tags Account
//
// @Accept json
// @Produce json
//
// @Param account body dtos.AccountCredentialsDto true "Account credentials"
//
// @Success 200 {object} dtos.AuthenticationResponseDataDto
// @Failure 401 {object} dtos.HTTPErrorDto
// @Failure 400 {object} dtos.HTTPErrorDto
//
// @Router /accounts/login [post]
func (controller AccountController) HandlePostLogin(ctx *gin.Context) {
	logger.RTrace("Service found. Getting yaml from it")

	var accountCredentials dtos.AccountCredentialsDto
	if err := ctx.ShouldBindJSON(&accountCredentials); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Get user information from the user management
	_, userCouldBeFound := controller.UserManagement.TryGetUserInformation(
		accountCredentials.Username,
		accountCredentials.Password,
	)

	// If the user could not be found, access should be not granted
	if !userCouldBeFound {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// Otherwise, return the role
	authData := dtos.AuthenticationResponseDataDto{
		Role:    "user",
		IsValid: true,
	}

	ctx.JSON(http.StatusOK, authData)
}
