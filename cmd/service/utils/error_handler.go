package utils

import (
	"OperatorAutomation/cmd/service/webserver/dtos"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewError(ctx *gin.Context, status int, err error) {

	errorMessage :=  err.Error()
	er := dtos.HTTPErrorDto{
		Code:    status,
		Message: errorMessage,
	}

	log.Error("Request produced an error: " + errorMessage)

	ctx.JSON(status, er)
}
