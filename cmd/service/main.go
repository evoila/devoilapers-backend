package main

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/webserver"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/kibana"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name: "start",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "configfile",
					Aliases: []string{"c"},
					Value:   "appconfig.json",
					Usage:   "Application configuration file. Includes port, certifcates, users etc...",
				},
			},
			Aliases: []string{"s"},
			Usage:   "Start webserver",
			Action: func(context *cli.Context) error {
				return initialize(context, false)
			},
		}, {
			Name: "demo",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "configfile",
					Aliases: []string{"c"},
					Value:   "appconfig.json",
					Usage:   "Application configuration file. Includes port, certifcates, users etc...",
				},
			},
			Aliases: []string{"s"},
			Usage:   "Start webserver demo",
			Action: func(context *cli.Context) error {
				return initialize(context, true)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Exit application")
	}
}

func initialize(c *cli.Context, demoMode bool) error {
	// Import config file
	filepath := c.String("configfile")
	parsedConfig, err := config.LoadConfigurationFromFile(filepath)
	if err != nil {
		log.Error("Config file in path could not be found or parsed. Ensure file exists and is valid json")
		log.Fatal(err)
		return err
	}

	//Apply loglevel
	ApplyGlobalLogConfigurations(parsedConfig)

	var appCore *core.Core

	if demoMode {
		log.Info("App launched in demo mode")
		appCore = InitializeDemoCore(parsedConfig)
	} else {
		// Create the core of the app
		appCore = InitializeCore(parsedConfig)
	}

	// Start webserver
	log.Info("Starting the webserver")
	err = webserver.StartWebserver(parsedConfig, appCore)
	if err != nil {
		log.Error("Webserver start failed")
		log.Fatal(err)
	}

	return err
}

// Create the core object that the service is interacting with
func InitializeCore(appconfig config.RawConfig) *core.Core {

	// TODO: Add concrete just like here service providers here
	var esp service.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.YamlTemplatePath,
	)

	var kb service.IServiceProvider = kibana.CreateKibanaProvider(
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.YamlTemplatePath,
	)

	var pg service.IServiceProvider = postgres.CreatePostgresProvider(
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.YamlTemplatePath,
		kubernetes.NginxInformation(appconfig.Kubernetes.Nginx),
	)

	return core.CreateCore([]*service.IServiceProvider{
		&esp,
		&kb,
		&pg,
	})
}
func InitializeDemoCore(appconfig config.RawConfig) *core.Core {
	var dp service.IServiceProvider = dummy.CreateDummyProvider()

	return core.CreateCore([]*service.IServiceProvider{
		&dp,
	})

}

// Set the loglevel from the config globally
func ApplyGlobalLogConfigurations(rawConfig config.RawConfig) {
	switch rawConfig.LogLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "warning":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	default:
		log.SetLevel(log.DebugLevel)
		log.Warn("Invalid log level found. Valid values are: trace, debug, warning, error. Fallback to debug level")
		break
	}
}
