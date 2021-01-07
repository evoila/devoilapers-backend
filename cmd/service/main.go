package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/cmd/service/controller"
	"OperatorAutomation/pkg/demolib"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toorop/gin-logrus"
	"github.com/swaggo/files"
)


// @title Operator Automation Backend API
// @version 1.0
// @description Operator Automation Backend API overview.
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
//
// @securityDefinitions.basic BasicAuth
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

	// Basic authentication users
	auth := gin.BasicAuth(gin.Accounts{
		"admin": "masterkey",
	})

	// Define routes
	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("/login", controller.HandlePostLogin)
		}
		servicestore := v1.Group("/servicestore", auth)
		{
			servicestore.GET("/info", controller.HandleGetServiceStoreOverview)
			servicestore.GET("/yaml/:servicetype", controller.HandleGetServiceStoreItemYaml)
		}
		services := v1.Group("/services", auth)
		{
			services.POST("/create/:servicetype", controller.HandlePostCreateServiceInstance)
			services.POST("/update/:serviceid", controller.HandlePostUpdateServiceInstance)
			services.DELETE("/:serviceid", controller.HandlePostDeleteServiceInstance)
			services.GET("/info", controller.HandleGetServiceInstanceDetailsForAllInstances)
			services.GET("/info/:serviceid", controller.HandleGetServiceInstanceDetails)
			services.GET("/yaml/:serviceid", controller.HandleGetServiceInstanceYaml)
		}
	}

	// Define swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server on 127.0.0.1:8080
	router.Run(":8080")
}
