package log

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoggerFactory defines the signature for functions that can create a new logger.
type LoggerFactory func(level Level, config *viper.Viper) (Logger, error)

// loggerFactories holds the registered logger factories.
var loggerFactories = make(map[string]LoggerFactory)

// RegisterFactory adds a logger factory to the registry.
// This function is called by logger implementations in their init() function.
func RegisterFactory(name string, factory LoggerFactory) {
	if factory == nil {
		panic("Logger factory " + name + " is nil")
	}
	if _, ok := loggerFactories[name]; ok {
		// Already registered, noop
		return
	}
	loggerFactories[name] = factory
}

// CreateLoggerFromConfig creates a logger instance based on the provided Viper configuration.
// It can create a single logger or a multi-logger if multiple outputs are configured.
func CreateLoggerFromConfig(v *viper.Viper) (Logger, error) {
	if v == nil {
		// Default to a simple console logger if no configuration is provided.
		return NewConsoleLogger(InfoLevel), nil
	}

	level := parseLevel(v.GetString("log.level"))

	loggerBackend := v.Sub("log.loggers")

	if loggerBackend == nil {
		// If no 'loggers' section, default to console logger
		return NewConsoleLogger(level), nil
	}

	allLoggers := loggerBackend.AllSettings()
	loggers := make([]Logger, 0, len(allLoggers))

	for key := range allLoggers {
		loggerConfig := loggerBackend.Sub(key)
		if !loggerConfig.GetBool("enabled") {
			continue
		}

		driver := loggerConfig.GetString("driver")
		factory, ok := loggerFactories[driver]
		if !ok {
			return nil, fmt.Errorf("logger driver %s not found", driver)
		}

		logger, err := factory(level, loggerConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger %s: %w", key, err)
		}
		loggers = append(loggers, logger)
	}

	// Shrink slice to fit the actual number of enabled loggers.
	if len(loggers) < cap(loggers) {
		_loggers := make([]Logger, len(loggers))
		copy(_loggers, loggers)
		loggers = _loggers
	}

	// Return the appropriate logger.
	switch len(loggers) {
	case 0:
		// If no loggers are enabled, default to a console logger.
		return NewConsoleLogger(level), nil
	case 1:
		return loggers[0], nil
	default:
		return NewMultiLogger(loggers...), nil
	}
}

// parseLevel converts a string log level to the Level type.
func parseLevel(levelStr string) Level {
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "panic":
		return PanicLevel
	default:
		return InfoLevel
	}
}
