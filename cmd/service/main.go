package main

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/pkg/elasticsearch"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	//host:= "https://127.0.0.1:49153"
	//tokenCrhis :=  "eyJhbGciOiJSUzI1NiIsImtpZCI6Il9aV0F1RnNOSEV0VllHVWt3UmVPTFlGTWpFb1g2RHRCNzA2TVRsV2NLRlkifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InVzZXIwLXRva2VuLWdiYmY4Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InVzZXIwIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiZjFhY2ExMjktMzlhNS00MzgwLTlhNWItN2E4Y2FkMzhhNTU3Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6dXNlcjAifQ.FhOq24sC1tRIwODSKTWUvkhGhzXqCl94hhmTwChpdH88SC_0csnyU9-EFiTHTKdYMnJrnxCxwxedyta_WESRAz62y9YPZ6FeqA2a9ttscRXlNHQmvb0CBy4W0KSb3_nZB8y8wXVA2FfMAyFU-RcQiT6PccTV_l9kNdjuqp7x_HRLWDw-QUvDpajiQZS-DVp5pB3bii49Gslhdm3Yp6iv4O8g8pgRoSZ2NC79aTNyCGpVh4NR54zRMMrP_9XXCIX26cdYtlJJAggPl3qtAM53acSRBtKHgIRZ4OaUgoRnDNZmiRdVF7fk7m1M00iJ2Rd55t0P8sUSAUJpeA3t_PlsAA"
	//_ = tokenCrhis
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IlB0ZnVSaW4xdUJTWjQxQzRTZE9YRmplcHI1Y0tibWxvS2hVT2d1NFlKQUEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InVzZXIwLXRva2VuLWdqMnhrIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InVzZXIwIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNzBiZTY4ZDMtMDgwZS00ODNlLWI3MTQtZDBjOTUxNTc3ZDJiIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6dXNlcjAifQ.fTr804ydneXofYaoMoW0V4OWyyYkO1dPn0MnQgb2cKXW8MP2qSOaI8xkiN5IzzshpGxwFGipD7G1f8ox9LgmzFFHFKG2TEw05QKKK9YMPctRih_sbuGAgTWcMOpK71ssQ6-afBOHjTospLhND0X30in0J-vcfLXGlAypBisfBuKyHKoUO2GHbgX2fW_3-xwsVQvITxXiWSDyDA-XGiVFgWL6ydfIWZ_iF0NI3qh8NYsit4BQmjaf1boMaImh_XTNAH4eBYmYkF4WIcQWoMfNCAQZNjRHC3CiXjfKIY4t38p7rXHp_uvnQ8nfNfYYVEG-Ci66Zi2iIDVLyKJjgjGc-w"
	//
	//obj, _ := cc.CreateKubernetesWrapper(host, token)
	//_ = obj
	//obj.Apply("apiVersion: elasticsearch.k8s.elastic.co/v1\nkind: Elasticsearch\nmetadata:\n  name: wibu\nspec:\n  version: 7.10.0\n  nodeSets:\n  - name: ganmo\n    count: 1\n    config:\n      node.store.allow_mmap: false")
	//log.SetLevel(log.TraceLevel)
	//log.SetLevel(log.TraceLevel)

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
				return initialize(context,false)
			},
		},{
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
				return initialize(context,true)
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
		err = StartWebserver(parsedConfig, appCore)
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
		appconfig.Kubernetes.CertificateAuthority)


	return core.CreateCore([]*service.IServiceProvider{
		&esp,
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
