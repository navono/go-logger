package main

import (
	"log"

	"github.com/navono/go-logger"
)

func main() {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.DebugLevel,
		ConsoleJSONFormat: true,
		EnableFile:        true,
		FileLevel:         logger.InfoLevel,
		FileJSONFormat:    true,
		Filename:          "log.log",
		FileMaxSize:       1,
		FileMaxAge:        1,
		Skip:              1,
	}
	zapLogger, err := logger.NewLogger(&config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	zapLogger.Infof("Starting with zap")

	contextLogger := zapLogger.WithFields(logger.Fields{"key1": "value1"})
	contextLogger.Infof("Zap is awesome")

	//c := logger.GetConcreteLogger(zapLogger)
	//if c != nil {
	//	sl := c.(*zap.SugaredLogger)
	//	zl := sl.Desugar()
	//	zl.Debug("concrete zap logger")
	//}

	logrusLogger, err := logger.NewLogger(&config, logger.InstanceLogrusLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	logrusLogger.Debugf("Starting with logrus")
	logrusLogger.Infof("Logrus is awesome")

	contextLogger = logrusLogger.WithFields(logger.Fields{"key1": "value1"})
	contextLogger.Debugf("Starting with context logrus")
	contextLogger.Infof("Logrus is awesome")
}
