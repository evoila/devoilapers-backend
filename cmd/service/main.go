package main

import (
	_ "OperatorAutomation/api" //Indirect use for swagger
	"OperatorAutomation/pkg/demolib"
	log "github.com/sirupsen/logrus"
)


func main() {
	log.SetLevel(log.TraceLevel)
	log.Info("Application started")

	// How to log
	log.Info(demolib.HelloWorld())

	// Start the webserver
	StartWebserver()
}
