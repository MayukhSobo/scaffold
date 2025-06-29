package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// ConsoleLogger implements Logger interface for console output
type ConsoleLogger struct {
	logger      zerolog.Logger
	level       LogLevel
	contextData map[string]any
	writer      io.Writer
}

// NewConsoleLogger creates a new console logger with specified level
func NewConsoleLogger(level LogLevel) Logger {
	return NewConsoleLoggerWithWriter(level, os.Stdout, true)
}

// NewConsoleLoggerWithWriter creates a console logger with custom writer and colorization
func NewConsoleLoggerWithWriter(level LogLevel, writer io.Writer, colorized bool) Logger {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(parseLogLevel(string(level)))

	var logger zerolog.Logger
	if colorized {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(writer).With().Timestamp().Caller().Logger()
	}

	return &ConsoleLogger{
		logger:      logger,
		level:       level,
		contextData: make(map[string]any),
		writer:      writer,
	}
}

// parseLogLevel converts string to zerolog level
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// addFields adds fields to the zerolog event
func (l *ConsoleLogger) addFields(event *zerolog.Event, fields []Field) *zerolog.Event {
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
func (l *ConsoleLogger) Debug(msg string, fields ...Field) {
	event := l.logger.Debug()
	l.addFields(event, fields).Msg(msg)
}

// Info logs an info message
func (l *ConsoleLogger) Info(msg string, fields ...Field) {
	event := l.logger.Info()
	l.addFields(event, fields).Msg(msg)
}

// Warn logs a warning message
func (l *ConsoleLogger) Warn(msg string, fields ...Field) {
	event := l.logger.Warn()
	l.addFields(event, fields).Msg(msg)
}

// Error logs an error message
func (l *ConsoleLogger) Error(msg string, fields ...Field) {
	event := l.logger.Error()
	l.addFields(event, fields).Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *ConsoleLogger) Fatal(msg string, fields ...Field) {
	event := l.logger.Fatal()
	l.addFields(event, fields).Msg(msg)
}

// Panic logs a panic message and panics
func (l *ConsoleLogger) Panic(msg string, fields ...Field) {
	event := l.logger.Panic()
	l.addFields(event, fields).Msg(msg)
}

// WithFields creates a new logger with additional context fields
func (l *ConsoleLogger) WithFields(fields ...Field) Logger {
	newContextData := make(map[string]any)

	// Copy existing context data
	for k, v := range l.contextData {
		newContextData[k] = v
	}

	// Add new fields
	for _, field := range fields {
		newContextData[field.Key] = field.Value
	}

	return &ConsoleLogger{
		logger:      l.logger,
		level:       l.level,
		contextData: newContextData,
		writer:      l.writer,
	}
}

// WithContext creates a new logger with context (for future use with request tracing)
func (l *ConsoleLogger) WithContext(ctx context.Context) Logger {
	// For now, just return a copy. This can be extended for request tracing
	return &ConsoleLogger{
		logger:      l.logger,
		level:       l.level,
		contextData: l.contextData,
		writer:      l.writer,
	}
}
