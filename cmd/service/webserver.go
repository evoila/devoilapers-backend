package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/controller"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
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
// @Schemes https
// @BasePath /api/v1
// @query.collection.format multi
//
// @x-extension-openapi {"example": "value on a json format"}
func StartWebserver(config config.RawConfig) error {
	log := logrus.New()
	// Set to global log level
	log.SetLevel(logrus.GetLevel())

	// Service Router
	router := gin.New()

	// Logging and recovery middleware
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// Allow cross origins
	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	c.AllowCredentials = true
	c.AddAllowHeaders("authorization")
	router.Use(cors.New(c))

	// Basic authentication users
	// Import them from the given config
	validAccounts := gin.Accounts{}
	for _, user := range config.User {
		validAccounts[user.Name] = user.Password
	}
	auth := gin.BasicAuth(validAccounts)

	// Define routes
	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("/login", func(context *gin.Context) {
				controller.HandlePostLogin(context, validAccounts)
			})
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
			services.POST("/action/:serviceid/:actioncommand", controller.HandlePostServiceInstanceAction)
			services.DELETE("/:serviceid", controller.HandleDeleteServiceInstance)
			services.GET("/info", controller.HandleGetServiceInstanceDetailsForAllInstances)
			services.GET("/info/:serviceid", controller.HandleGetServiceInstanceDetails)
			services.GET("/yaml/:serviceid", controller.HandleGetServiceInstanceYaml)
		}
	}

	// Define swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Debugf("Visit https://127.0.0.1:%d/swagger/index.html to see the swagger document", config.Port)

	// Start server
	return router.RunTLS(
		":"+strconv.Itoa(config.Port),
		config.WebserverSllCertificate.PublicKeyFilePath,
		config.WebserverSllCertificate.PrivateKeyFilePath)
}
