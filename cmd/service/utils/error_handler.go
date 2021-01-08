package utils

import (
	"OperatorAutomation/cmd/service/dtos"
	"github.com/gin-gonic/gin"
)

func NewError(ctx *gin.Context, status int, err error) {
	er := dtos.HTTPErrorDto{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
