package log

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileLoggerConfig contains configuration for file logging with rotation
type FileLoggerConfig struct {
	Filename   string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool // compress rotated files
	JsonFormat bool // use JSON format
}

// FileLogger implements Logger interface for file output with rotation
type FileLogger struct {
	logger      zerolog.Logger
	level       LogLevel
	contextData map[string]any
	lumberjack  *lumberjack.Logger
	config      *FileLoggerConfig
}

// NewFileLogger creates a new file logger with rotation
func NewFileLogger(level LogLevel, config *FileLoggerConfig) Logger {
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

// addFields adds fields to the zerolog event
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

// Debug logs a debug message
func (l *FileLogger) Debug(msg string, fields ...Field) {
	event := l.logger.Debug()
	l.addFields(event, fields).Msg(msg)
}

// Info logs an info message
func (l *FileLogger) Info(msg string, fields ...Field) {
	event := l.logger.Info()
	l.addFields(event, fields).Msg(msg)
}

// Warn logs a warning message
func (l *FileLogger) Warn(msg string, fields ...Field) {
	event := l.logger.Warn()
	l.addFields(event, fields).Msg(msg)
}

// Error logs an error message
func (l *FileLogger) Error(msg string, fields ...Field) {
	event := l.logger.Error()
	l.addFields(event, fields).Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *FileLogger) Fatal(msg string, fields ...Field) {
	event := l.logger.Fatal()
	l.addFields(event, fields).Msg(msg)
}

// Panic logs a panic message and panics
func (l *FileLogger) Panic(msg string, fields ...Field) {
	event := l.logger.Panic()
	l.addFields(event, fields).Msg(msg)
}

// WithFields creates a new logger with additional context fields
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

// WithContext creates a new logger with context
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

// Close closes the file logger and flushes any remaining logs
func (l *FileLogger) Close() error {
	return l.lumberjack.Close()
}
