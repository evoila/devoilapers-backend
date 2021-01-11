package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/controller"
	"OperatorAutomation/pkg/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
	"strconv"
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
func StartWebserver(appconfig config.RawConfig, core *core.Core) error {
	log := logrus.New()
	// Set to global log level
	log.SetLevel(logrus.GetLevel())

	// Service Router
	router := gin.New()

	// Logging and recovery middleware
	router.Use(ginlogrus.Logger(log), gin.Recovery())
	// Allow cross origins
	router.Use(cors.Default())

	// Basic authentication users
	// Import them from the given config
	validAccounts := gin.Accounts{}
	for _, user := range core.UserContextManagement.Users {
		validAccounts[user.GetName()] = user.GetPassword()
	}
	auth := gin.BasicAuth(validAccounts)

	// Define routes
	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("/login", func(context *gin.Context) {
				controller.HandlePostLogin(context, core)
			})
		}
		servicestore := v1.Group("/servicestore", auth)
		{
			servicestore.GET("/info", func(context *gin.Context) {
				controller.HandleGetServiceStoreOverview(context, core)
			})
			servicestore.GET("/yaml/:servicetype", func(context *gin.Context) {
				controller.HandleGetServiceStoreItemYaml(context, core)
			})
		}
		services := v1.Group("/services", auth)
		{
			services.POST("/create/:servicetype", func(context *gin.Context) {
				controller.HandlePostCreateServiceInstance(context, core)
			})
			services.POST("/update/:serviceid", func(context *gin.Context) {
				controller.HandlePostUpdateServiceInstance(context, core)
			})
			services.POST("/action/:serviceid/:actioncommand", func(context *gin.Context) {
				controller.HandlePostServiceInstanceAction(context, core)
			})
			services.DELETE("/:serviceid", func(context *gin.Context) {
				controller.HandleDeleteServiceInstance(context, core)
			})
			services.GET("/info", func(context *gin.Context) {
				controller.HandleGetServiceInstanceDetailsForAllInstances(context, core)
			})
			services.GET("/info/:serviceid", func(context *gin.Context) {
				controller.HandleGetServiceInstanceDetails(context, core)
			})
			services.GET("/yaml/:serviceid", func(context *gin.Context) {
				controller.HandleGetServiceInstanceYaml(context, core)
			})
		}
	}

	// Define swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Debugf("Visit https://127.0.0.1:%d/swagger/index.html to see the swagger document", appconfig.Port)

	// Start server
	return router.RunTLS(
		":"+strconv.Itoa(appconfig.Port),
		appconfig.WebserverSllCertificate.PublicKeyFilePath,
		appconfig.WebserverSllCertificate.PrivateKeyFilePath)
}
