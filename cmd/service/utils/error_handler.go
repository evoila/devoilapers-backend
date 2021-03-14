package utils

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"github.com/gin-gonic/gin"
)

func NewError(ctx *gin.Context, status int, err error) {
	er := dtos.HTTPErrorDto{
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
