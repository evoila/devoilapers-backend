package main

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/service"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.SetLevel(log.TraceLevel)

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
			Action: func(c *cli.Context) error {
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

				// Create the core of the app
				appCore := InitializeCore(parsedConfig)

				// Start webserver
				log.Info("Starting the webserver")
				err = StartWebserver(parsedConfig, appCore)
				if err != nil {
					log.Error("Webserver start failed")
					log.Fatal(err)
				}

				return err
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

// Create the core object that the service is interacting with
func InitializeCore(appconfig config.RawConfig) *core.Core {
	// TODO: Add concrete service providers here
	return core.CreateCore([]*service.IServiceProvider{})
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
