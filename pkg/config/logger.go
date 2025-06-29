package config

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/MayukhSobo/scaffold/pkg/log"

	"github.com/spf13/viper"
)

// CreateLoggerFromConfig creates a logger based on viper configuration
func CreateLoggerFromConfig(v *viper.Viper) log.Logger {
	if v == nil {
		// Return default logger if no config
		return log.NewConsoleLogger(log.InfoLevel)
	}

	// Parse log level
	level := parseLevel(v.GetString("log.level"))

	// Collect enabled loggers using extensible pattern
	var loggers []log.Logger

	// Add console logger if enabled
	if v.GetBool("log.console_logger.enabled") {
		consoleLogger := createConsoleLogger(level, v)
		loggers = append(loggers, consoleLogger)
	}

	// Add file logger if enabled
	if v.GetBool("log.file_logger.enabled") && v.GetString("log.file_logger.filename") != "" {
		fileLogger := createFileLogger(level, v)
		loggers = append(loggers, fileLogger)
	}

	// Future loggers can be added here when implementations exist:
	// if v.GetBool("log.datadog_logger.enabled") {
	//     datadogLogger := createDatadogLogger(level, v)
	//     loggers = append(loggers, datadogLogger)
	// }

	// Return appropriate logger based on number of enabled loggers
	switch len(loggers) {
	case 0:
		// Default to console logger if no loggers configured
		return log.NewConsoleLogger(level)
	case 1:
		return loggers[0]
	default:
		return log.NewMultiLogger(loggers...)
	}
}

// createConsoleLogger creates a console logger with specific configuration
func createConsoleLogger(level log.LogLevel, v *viper.Viper) log.Logger {
	if v.GetBool("log.console_logger.json_format") {
		// For JSON format, disable colors
		return log.NewConsoleLoggerWithWriter(level, nil, false)
	}
	return log.NewConsoleLogger(level) // Uses colors based on config
}

// createFileLogger creates a file logger with specific configuration
func createFileLogger(level log.LogLevel, v *viper.Viper) log.Logger {
	// Get configuration values
	directory := v.GetString("log.file_logger.directory")
	filename := v.GetString("log.file_logger.filename")

	// Resolve the full file path
	fullPath := resolveLogFilePath(directory, filename)

	// Ensure the directory exists
	if err := ensureLogDirectory(filepath.Dir(fullPath)); err != nil {
		// Log the error and fall back to console logger
		fmt.Printf("Warning: %v, falling back to console logger\n", err)
		return log.NewConsoleLogger(level)
	}

	fileConfig := &log.FileLoggerConfig{
		Filename:   fullPath,
		MaxSize:    v.GetInt("log.file_logger.max_size"),
		MaxBackups: v.GetInt("log.file_logger.max_backups"),
		MaxAge:     v.GetInt("log.file_logger.max_age"),
		Compress:   v.GetBool("log.file_logger.compress"),
		JsonFormat: v.GetBool("log.file_logger.json_format"),
	}
	return log.NewFileLogger(level, fileConfig)
}

// resolveLogFilePath creates the full path for the log file
func resolveLogFilePath(directory, filename string) string {
	// If directory is empty, use default logs directory
	if directory == "" {
		directory = "logs"
	}

	// If directory is not absolute, make it relative to current working directory
	if !filepath.IsAbs(directory) {
		if cwd, err := os.Getwd(); err == nil {
			directory = filepath.Join(cwd, directory)
		}
	}

	return filepath.Join(directory, filename)
}

// ensureLogDirectory creates the log directory if it doesn't exist
func ensureLogDirectory(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}
	return nil
}

// parseLevel converts string level to log.LogLevel
func parseLevel(levelStr string) log.LogLevel {
	switch levelStr {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}
