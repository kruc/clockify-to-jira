package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func configureLogger() {
	log.SetFormatter(&log.TextFormatter{})
	if globalConfig.logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)
	if globalConfig.logOutput != "stdout" {
		logFile, _ := os.OpenFile(globalConfig.logOutput, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(logFile)
	}
	log.SetLevel(log.InfoLevel)
}
