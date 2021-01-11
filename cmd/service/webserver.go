package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/controller"
	"OperatorAutomation/cmd/service/user"
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
	for _, user := range appconfig.Users {
		validAccounts[user.GetName()] = user.GetPassword()
	}
	auth := gin.BasicAuth(validAccounts)

	// Define controller with access to the core component
	baseController := controller.BaseController{Core: core, UserManagement: user.CreateUserManagement(appconfig.Users)}
	serviceStoreController := controller.ServiceStoreController{BaseController: baseController}
	serviceController := controller.ServiceController{BaseController: baseController}
	accountController := controller.AccountController{BaseController: baseController}

	// Define routes
	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("/login", accountController.HandlePostLogin)
		}
		servicestore := v1.Group("/servicestore", auth)
		{
			servicestore.GET("/info", serviceStoreController.HandleGetServiceStoreOverview)
			servicestore.GET("/yaml/:servicetype", serviceStoreController.HandleGetServiceStoreItemYaml)
		}
		services := v1.Group("/services", auth)
		{
			services.POST("/create/:servicetype", serviceController.HandlePostCreateServiceInstance)
			services.POST("/update/:serviceid", serviceController.HandlePostUpdateServiceInstance)
			services.POST("/action/:serviceid/:actioncommand", serviceController.HandlePostServiceInstanceAction)
			services.DELETE("/:serviceid", serviceController.HandleDeleteServiceInstance)
			services.GET("/info", serviceController.HandleGetServiceInstanceDetailsForAllInstances)
			services.GET("/info/:serviceid", serviceController.HandleGetServiceInstanceDetails)
			services.GET("/yaml/:serviceid", serviceController.HandleGetServiceInstanceYaml)
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
