package log

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/MayukhSobo/scaffold/pkg/utils"
)

// FileLoggerConfig contains configuration for file logging with rotation.
type FileLoggerConfig struct {
	Filename   string `mapstructure:"filename"`
	Directory  string `mapstructure:"directory"`
	MaxSize    int    `mapstructure:"max_size"`    // megabytes
	MaxBackups int    `mapstructure:"max_backups"` // number of backups
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`    // compress rotated files
	JsonFormat bool   `mapstructure:"json_format"` // use JSON format
}

// FileLogger implements Logger interface for file output with rotation.
type FileLogger struct {
	logger      zerolog.Logger
	level       Level
	contextData map[string]any
	lumberjack  *lumberjack.Logger
	config      *FileLoggerConfig
}

func init() {
	RegisterFactory("file", NewFileLoggerFromConfig)
}

// NewFileLoggerFromConfig creates a new file logger from a Viper configuration.
func NewFileLoggerFromConfig(level Level, v *viper.Viper) (Logger, error) {
	var config FileLoggerConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	fullPath := utils.ResolveLogFilePath(config.Directory, config.Filename)
	if err := utils.EnsureLogDirectory(filepath.Dir(fullPath)); err != nil {
		return nil, err
	}

	// The existing NewFileLogger expects a config with the full path.
	fileLoggerConfig := &FileLoggerConfig{
		Filename:   fullPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		JsonFormat: config.JsonFormat,
	}

	return NewFileLogger(level, fileLoggerConfig), nil
}

// NewFileLogger creates a new file logger with rotation.
func NewFileLogger(level Level, config *FileLoggerConfig) Logger {
	// Set defaults if not provided
	if config.MaxSize == 0 {
		config.MaxSize = 100 // 100MB
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 3
	}
	if config.MaxAge == 0 {
		config.MaxAge = 7 // 7 days
	}

	lj := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(parseLogLevel(string(level)))

	var logger zerolog.Logger
	if config.JsonFormat {
		logger = zerolog.New(lj).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: lj, NoColor: true}).With().Timestamp().Caller().Logger()
	}

	return &FileLogger{
		logger:      logger,
		level:       level,
		contextData: make(map[string]any),
		lumberjack:  lj,
		config:      config,
	}
}

// addFields adds fields to the zerolog event.
func (l *FileLogger) addFields(event *zerolog.Event, fields []Field) *zerolog.Event {
	// Add context data first
	for k, v := range l.contextData {
		event = event.Interface(k, v)
	}

	// Add provided fields
	for _, field := range fields {
		event = event.Interface(field.Key, field.Value)
	}
	return event
}

// Debug logs a debug message.
func (l *FileLogger) Debug(msg string, fields ...Field) {
	event := l.logger.Debug()
	l.addFields(event, fields).Msg(msg)
}

// Info logs an info message.
func (l *FileLogger) Info(msg string, fields ...Field) {
	event := l.logger.Info()
	l.addFields(event, fields).Msg(msg)
}

// Warn logs a warning message.
func (l *FileLogger) Warn(msg string, fields ...Field) {
	event := l.logger.Warn()
	l.addFields(event, fields).Msg(msg)
}

// Error logs an error message.
func (l *FileLogger) Error(msg string, fields ...Field) {
	event := l.logger.Error()
	l.addFields(event, fields).Msg(msg)
}

// Fatal logs a fatal message and exits.
func (l *FileLogger) Fatal(msg string, fields ...Field) {
	event := l.logger.Fatal()
	l.addFields(event, fields).Msg(msg)
}

// Panic logs a panic message and panics.
func (l *FileLogger) Panic(msg string, fields ...Field) {
	event := l.logger.Panic()
	l.addFields(event, fields).Msg(msg)
}

// Formatted logging methods
func (l *FileLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprintf(format, args...))
}

func (l *FileLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *FileLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprintf(format, args...))
}

func (l *FileLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprintf(format, args...))
}

func (l *FileLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l *FileLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprintf(format, args...))
}

// WithFields creates a new logger with additional context fields.
func (l *FileLogger) WithFields(fields ...Field) Logger {
	newContextData := make(map[string]any)

	// Copy existing context data
	for k, v := range l.contextData {
		newContextData[k] = v
	}

	// Add new fields
	for _, field := range fields {
		newContextData[field.Key] = field.Value
	}

	return &FileLogger{
		logger:      l.logger,
		level:       l.level,
		contextData: newContextData,
		lumberjack:  l.lumberjack,
		config:      l.config,
	}
}

// WithContext creates a new logger with context.
func (l *FileLogger) WithContext(ctx context.Context) Logger {
	// For now, just return a copy. This can be extended for request tracing
	return &FileLogger{
		logger:      l.logger,
		level:       l.level,
		contextData: l.contextData,
		lumberjack:  l.lumberjack,
		config:      l.config,
	}
}

// Close closes the file logger and flushes any remaining logs.
func (l *FileLogger) Close() error {
	return l.lumberjack.Close()
}
