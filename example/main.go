package main

import (
	"log"

	"github.com/navono/go-logger"
)

func main() {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: true,
		EnableFile:        true,
		FileLevel:         logger.Info,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}
	err := logger.NewLogger(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
	
  contextLogger := logger.WithFields(logger.Fields{"key1": "value1"})
  contextLogger.Debugf("Starting with zap")
  contextLogger.Infof("Zap is awesome")
  
  err = logger.NewLogger(config, logger.InstanceLogrusLogger)
    if err != nil {
      log.Fatalf("Could not instantiate log %s", err.Error())
    }
  contextLogger = logger.WithFields(logger.Fields{"key1": "value1"})
  contextLogger.Debugf("Starting with logrus")
  
  contextLogger.Infof("Logrus is awesome")
}
