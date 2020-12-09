package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/pkg/demolib"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/toorop/gin-logrus"
)


// @title Operator Automation Backend API
// @version 1.0
// @description Operator Automation Backend API overview.
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
//
// @BasePath /api/v1
// @query.collection.format multi
//
// @x-extension-openapi {"example": "value on a json format"}
func main() {
	log := logrus.New()

	// How to log
	log.Info(demolib.HelloWorld())

	// Service Router
	router := gin.New()

	// Logging and recovery middleware
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// Define routes
	v1 := router.Group("/api/v1")
	{
		root := v1.Group("/")
		{
			root.GET("", ShowHello)
		}
	}

	// Define swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server on 127.0.0.1:8080
	router.Run(":8080")
}
