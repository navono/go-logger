package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = "debug"

	// InfoLevel is the default logging priority.
	InfoLevel = "info"

	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = "warn"

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = "fatal"
)

const (
	// InstanceZapLogger will be used to create Zap instance for the logger
	InstanceZapLogger int = iota
	// InstanceLogrusLogger will be used to create Logrus instance for the logger
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// Configuration stores the config for the Logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string

	// EnableFile that whether to log information to a file
	EnableFile bool

	// FileLevel is the log need json format for file
	FileJSONFormat bool

	// FileLevel is the log level that for file
	FileLevel string

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string

	// FileMaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	FileMaxSize int

	// FileMaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	FileMaxAge int

	// Skip count for caller stack. Zap and logrus have different skip count
	Skip int

	zapLogger    *zapLogger
	logrusLogger *logrusLogger
}

// NewLogger returns an instance of Logger
func NewLogger(config *Configuration, loggerInstance int) (Logger, error) {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(*config)
		if err != nil {
			return nil, err
		}

		zl, ok := logger.(*zapLogger)
		if ok {
			config.zapLogger = zl
		}

		return logger, nil

	case InstanceLogrusLogger:
		logger, err := newLogrusLogger(*config)
		if err != nil {
			return nil, err
		}

		ll, ok := logger.(*logrusLogger)
		if ok {
			config.logrusLogger = ll
		}

		return logger, nil

	default:
		return nil, errInvalidLoggerInstance
	}
}

// GetConcreteLogger returns the underlying log instance
func GetConcreteLogger(log Logger) interface{} {
	switch l := log.(type) {
	case *zapLogger:
		return l.sugaredLogger
	case *logrusLogger:
		return l.logger
	}

	return nil
}

func (c *Configuration) SetLevel(l string) error {
	if c.zapLogger == nil && c.logrusLogger == nil {
		return fmt.Errorf("log not initialized")
	}

	if c.zapLogger != nil {
		c.zapLogger.atom.SetLevel(getZapLevel(l))
	}

	if c.logrusLogger != nil {
		if l, err := logrus.ParseLevel(l); err == nil {
			c.logrusLogger.logger.SetLevel(l)
		}
	}

	return nil
}
