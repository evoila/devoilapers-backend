package main

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/cmd/service/webserver"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/kibana"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres"
	"OperatorAutomation/pkg/utils/logger"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"math/rand"
	"net/url"
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
		logger.RInfo("Exit application")
	}
}

func initialize(c *cli.Context, demoMode bool) error {
	// Import config file
	filepath := c.String("configfile")
	parsedConfig, err := config.LoadConfigurationFromFile(filepath)
	if err != nil {
		logger.RError(err,"Config file in path could not be found or parsed. Ensure file exists and is valid json")
		log.Fatal(err)
		return err
	}

	//Apply loglevel
	ApplyGlobalLogConfigurations(parsedConfig)

	var appCore *core.Core

	if demoMode {
		log.Info("App launched in demo mode.")
		appCore = InitializeDemoCore(parsedConfig)
	} else {
		// Create the core of the app
		appCore = InitializeCore(parsedConfig)
	}

	// Start webserver
	logger.RInfo("Starting the webserver.")
	err = webserver.StartWebserver(parsedConfig, appCore)
	if err != nil {
		logger.RError(err,"Webserver start failed.")
		log.Fatal(err)
	}

	return err
}

// Create the core object that the service is interacting with
func InitializeCore(appconfig config.RawConfig) *core.Core {
	url, err := url.Parse(appconfig.Kubernetes.Server)
	if err != nil {
		logger.RError(err, "Could not parse kubernetes server url: " + appconfig.Kubernetes.Server)
		panic(err)
	}

	hostname := url.Hostname()

	// TODO: Add concrete just like here service providers here
	var esp provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		hostname,
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.ResourcesTemplatesPath,
	)

	var kb provider.IServiceProvider = kibana.CreateKibanaProvider(
		hostname,
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.ResourcesTemplatesPath,
	)


	var pg provider.IServiceProvider = postgres.CreatePostgresProvider(
		hostname,
		appconfig.Kubernetes.Server,
		appconfig.Kubernetes.CertificateAuthority,
		appconfig.Kubernetes.Operators.Postgres.PgoUrl,
		appconfig.Kubernetes.Operators.Postgres.PgoVersion,
		appconfig.Kubernetes.Operators.Postgres.PgoCaPath,
		appconfig.Kubernetes.Operators.Postgres.PgoUsername,
		appconfig.Kubernetes.Operators.Postgres.PgoPassword,
		appconfig.ResourcesTemplatesPath,
		kubernetes.NginxInformation(appconfig.Kubernetes.Nginx),
	)

	return core.CreateCore([]*provider.IServiceProvider{
		&esp,
		&kb,
		&pg,
	})
}
func InitializeDemoCore(appconfig config.RawConfig) *core.Core {
	var dp provider.IServiceProvider = dummy.CreateDummyProvider()

	return core.CreateCore([]*provider.IServiceProvider{
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
		logger.RWarn("Invalid log level found. Valid values are: trace, debug, warning, error. Fallback to debug level")
		break
	}
}
