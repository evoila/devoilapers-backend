package controller

import (
	"OperatorAutomation/cmd/service/dtos"
	"OperatorAutomation/cmd/service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
func HandlePostLogin(ctx *gin.Context, accounts gin.Accounts) {
	var accountCredentials dtos.AccountCredentialsDto

	if err := ctx.ShouldBindJSON(&accountCredentials); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	password, userCouldBeFound := accounts[accountCredentials.Username]
	if !userCouldBeFound || password != accountCredentials.Password {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	authData := dtos.AuthenticationResponseDataDto{
		Role:    "admin",
		IsValid: true,
	}

	ctx.JSON(http.StatusOK, authData)
}
