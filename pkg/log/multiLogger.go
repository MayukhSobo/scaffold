package log

import (
	"context"
)

// MultiLogger implements Logger interface and forwards logs to multiple loggers
// This allows combining console and file logging or any other logger implementations
type MultiLogger struct {
	loggers     []Logger
	contextData map[string]any
}

// NewMultiLogger creates a new multi-logger that forwards to multiple logger implementations
func NewMultiLogger(loggers ...Logger) Logger {
	return &MultiLogger{
		loggers:     loggers,
		contextData: make(map[string]any),
	}
}

// Debug logs a debug message to all underlying loggers
func (m *MultiLogger) Debug(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Debug(msg, fields...)
	}
}

// Info logs an info message to all underlying loggers
func (m *MultiLogger) Info(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Info(msg, fields...)
	}
}

// Warn logs a warning message to all underlying loggers
func (m *MultiLogger) Warn(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Warn(msg, fields...)
	}
}

// Error logs an error message to all underlying loggers
func (m *MultiLogger) Error(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Error(msg, fields...)
	}
}

// Fatal logs a fatal message to all underlying loggers and exits
func (m *MultiLogger) Fatal(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Fatal(msg, fields...)
	}
}

// Panic logs a panic message to all underlying loggers and panics
func (m *MultiLogger) Panic(msg string, fields ...Field) {
	for _, logger := range m.loggers {
		logger.Panic(msg, fields...)
	}
}

// WithFields creates a new multi-logger with additional context fields
func (m *MultiLogger) WithFields(fields ...Field) Logger {
	newLoggers := make([]Logger, len(m.loggers))
	for i, logger := range m.loggers {
		newLoggers[i] = logger.WithFields(fields...)
	}

	return &MultiLogger{
		loggers:     newLoggers,
		contextData: m.contextData,
	}
}

// WithContext creates a new multi-logger with context
func (m *MultiLogger) WithContext(ctx context.Context) Logger {
	newLoggers := make([]Logger, len(m.loggers))
	for i, logger := range m.loggers {
		newLoggers[i] = logger.WithContext(ctx)
	}

	return &MultiLogger{
		loggers:     newLoggers,
		contextData: m.contextData,
	}
}
