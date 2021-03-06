package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func configureLogger(logFormat, logOutput string) {
	log.SetFormatter(&log.TextFormatter{})
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)
	if logOutput != "stdout" {
		logFile, _ := os.OpenFile(logOutput, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(logFile)
	}
	log.SetLevel(log.InfoLevel)
}
