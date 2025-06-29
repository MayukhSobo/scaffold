package log

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"
)

func TestConsole(t *testing.T) {
	logger := NewConsoleLogger(InfoLevel)
	if logger == nil {
		t.Fatal("Console logger should not be nil")
	}

	// Test basic logging
	logger.Info("Test info message")
	logger.Debug("Test debug message", String("key", "value"))
	logger.Warn("Test warning message")
	logger.Error("Test error message", String("error", "test error"))
}

func TestConsoleWithWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewConsoleLoggerWithWriter(DebugLevel, &buf, false)

	logger.Info("Test message")

	if buf.Len() == 0 {
		t.Error("Nothing was written to buffer")
	}
}

func TestFile(t *testing.T) {
	logFile := "test_file.log"
	defer os.Remove(logFile) // Clean up

	config := &FileLoggerConfig{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 2,
		MaxAge:     1,
		Compress:   false,
		JsonFormat: true,
	}

	logger := NewFileLogger(InfoLevel, config)
	if logger == nil {
		t.Fatal("File logger should not be nil")
	}

	// Test logging
	logger.Info("File test message")
	logger.Error("File error", String("component", "test"))

	// Verify file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestMulti(t *testing.T) {
	var consoleBuf bytes.Buffer
	consoleLogger := NewConsoleLoggerWithWriter(InfoLevel, &consoleBuf, false)

	logFile := "test_multi.log"
	defer os.Remove(logFile)

	fileLogger := NewFileLogger(InfoLevel, &FileLoggerConfig{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
		JsonFormat: true,
	})

	multiLogger := NewMultiLogger(consoleLogger, fileLogger)
	if multiLogger == nil {
		t.Fatal("Multi logger should not be nil")
	}

	// Test logging to multiple outputs
	multiLogger.Info("Multi logger test message")

	// Verify console output
	if consoleBuf.Len() == 0 {
		t.Error("Nothing was written to console buffer")
	}

	// Verify file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created by multi logger")
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewConsoleLoggerWithWriter(DebugLevel, &buf, false)

	// Create logger with persistent fields
	contextLogger := logger.WithFields(
		String("request_id", "test-123"),
		String("user_id", "user-456"),
	)

	// Test that fields persist
	contextLogger.Info("Test message with context")
	contextLogger.Debug("Debug with context", String("additional", "data"))

	output := buf.String()
	if output == "" {
		t.Error("No output generated")
	}
}

func TestWithContext(t *testing.T) {
	logger := NewConsoleLogger(InfoLevel)
	ctx := context.Background()

	contextLogger := logger.WithContext(ctx)
	if contextLogger == nil {
		t.Fatal("Context logger should not be nil")
	}

	contextLogger.Info("Test message with context")
}

func TestFieldHelpers(t *testing.T) {
	var buf bytes.Buffer
	logger := NewConsoleLoggerWithWriter(InfoLevel, &buf, false)

	// Test different field types
	testTime := time.Now()
	testDuration := 5 * time.Second

	logger.Info("Testing field helpers",
		String("string_field", "test"),
		Int("int_field", 42),
		Int64("int64_field", int64(123)),
		Float64("float_field", 3.14),
		Bool("bool_field", true),
		Time("time_field", testTime),
		Duration("duration_field", testDuration),
		Any("any_field", map[string]string{"key": "value"}),
	)

	// Test error field
	testErr := os.ErrNotExist
	logger.Error("Test error logging", Error(testErr))

	output := buf.String()
	if output == "" {
		t.Error("No output generated for field helpers")
	}
}

func TestLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := NewConsoleLoggerWithWriter(DebugLevel, &buf, false)

	// Test all log levels
	logger.Debug("Debug level test")
	logger.Info("Info level test")
	logger.Warn("Warn level test")
	logger.Error("Error level test")

	output := buf.String()
	if output == "" {
		t.Error("No output generated for log levels")
	}

	// Note: Not testing Fatal and Panic as they would exit/panic the test
}

func TestFileConfig(t *testing.T) {
	config := &FileLoggerConfig{
		Filename:   "test_config.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
		JsonFormat: false,
	}

	defer os.Remove(config.Filename)

	logger := NewFileLogger(WarnLevel, config)
	if logger == nil {
		t.Fatal("File logger with config should not be nil")
	}

	logger.Warn("Config test message")

	// Verify file was created
	if _, err := os.Stat(config.Filename); os.IsNotExist(err) {
		t.Error("Log file was not created with custom config")
	}
}

func TestInterface(t *testing.T) {
	// Test that all our loggers implement the Logger interface
	var logger Logger

	// Test console logger
	logger = NewConsoleLogger(InfoLevel)
	testLoggerInterface(t, logger, "console")

	// Test file logger
	logFile := "interface_test.log"
	defer os.Remove(logFile)

	logger = NewFileLogger(InfoLevel, &FileLoggerConfig{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
		JsonFormat: true,
	})
	testLoggerInterface(t, logger, "file")

	// Test multi logger
	var buf bytes.Buffer
	consoleLogger := NewConsoleLoggerWithWriter(InfoLevel, &buf, false)
	multiLogger := NewMultiLogger(consoleLogger, logger)
	testLoggerInterface(t, multiLogger, "multi")
}

func testLoggerInterface(t *testing.T, logger Logger, loggerType string) {
	// Test all interface methods
	logger.Info("Info test", String("type", loggerType))
	logger.Warn("Warn test", String("type", loggerType))
	logger.Error("Error test", String("type", loggerType))

	// Test WithFields
	contextLogger := logger.WithFields(String("type", loggerType))
	contextLogger.Info("Fields test")

	// Test WithContext
	ctx := context.Background()
	ctxLogger := logger.WithContext(ctx)
	ctxLogger.Info("Context test")
}

// Benchmarks
func BenchmarkConsole(b *testing.B) {
	logger := NewConsoleLoggerWithWriter(InfoLevel, io.Discard, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message")
		}
	})
}

func BenchmarkStructured(b *testing.B) {
	logger := NewConsoleLoggerWithWriter(InfoLevel, io.Discard, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Structured benchmark message",
				String("key1", "value1"),
				Int("key2", 42),
				Bool("key3", true),
			)
		}
	})
}

func BenchmarkMulti(b *testing.B) {
	console := NewConsoleLoggerWithWriter(InfoLevel, io.Discard, false)
	file := NewFileLogger(InfoLevel, &FileLoggerConfig{
		Filename:   "/dev/null",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		JsonFormat: false,
	})
	logger := NewMultiLogger(console, file)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Multi benchmark message")
		}
	})
}

func TestParseLogLevel(t *testing.T) {
	testCases := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
		{"panic", PanicLevel},
		{"invalid", InfoLevel},
		{"", InfoLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := parseLevel(tc.input)
			if result != tc.expected {
				t.Errorf("Expected level %s, but got %s", tc.expected, result)
			}
		})
	}
}

func TestFileLoggerDebug(t *testing.T) {
	logFile := "test_file_debug.log"
	defer os.Remove(logFile) // Clean up

	config := &FileLoggerConfig{
		Filename:   logFile,
		JsonFormat: true,
	}

	logger := NewFileLogger(DebugLevel, config)
	logger.Debug("this should be logged")
	logger.Info("this should also be logged")

	// Verify file was created and has content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal("Could not read log file")
	}
	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}

func TestFileLoggerClose(t *testing.T) {
	logFile := "test_file_close.log"
	defer os.Remove(logFile)

	config := &FileLoggerConfig{Filename: logFile}
	logger := NewFileLogger(InfoLevel, config)

	// Cast to FileLogger to access Close method
	fileLogger, ok := logger.(*FileLogger)
	if !ok {
		t.Fatal("Could not cast to *FileLogger")
	}

	err := fileLogger.Close()
	if err != nil {
		t.Errorf("Error closing file logger: %v", err)
	}
}

func TestMultiLoggerDebug(t *testing.T) {
	var consoleBuf bytes.Buffer
	consoleLogger := NewConsoleLoggerWithWriter(DebugLevel, &consoleBuf, false)

	logFile := "test_multi_debug.log"
	defer os.Remove(logFile)
	fileLogger := NewFileLogger(DebugLevel, &FileLoggerConfig{Filename: logFile})

	multiLogger := NewMultiLogger(consoleLogger, fileLogger)
	multiLogger.Debug("multi-logger debug message")

	// Verify console output
	if consoleBuf.Len() == 0 {
		t.Error("Multi-logger did not write to console")
	}

	// Verify file output
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Multi-logger did not create log file")
	}
}

func TestLoggerWithEmptyFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewConsoleLoggerWithWriter(InfoLevel, &buf, false)

	// Test that no panic occurs with empty fields
	logger.WithFields().Info("Test with empty fields")
	logger.Info("Test with no fields")

	if buf.Len() == 0 {
		t.Error("No output with empty fields")
	}
}

func TestNewFileLoggerDefaults(t *testing.T) {
	logFile := "test_defaults.log"
	defer os.Remove(logFile)

	// Test with empty config to check defaults
	config := &FileLoggerConfig{Filename: logFile}
	logger := NewFileLogger(InfoLevel, config)

	// Cast to get access to internal config
	fileLogger, ok := logger.(*FileLogger)
	if !ok {
		t.Fatal("Could not cast to *FileLogger")
	}

	if fileLogger.config.MaxSize != 100 {
		t.Errorf("Expected default MaxSize=100, got %d", fileLogger.config.MaxSize)
	}
	if fileLogger.config.MaxBackups != 3 {
		t.Errorf("Expected default MaxBackups=3, got %d", fileLogger.config.MaxBackups)
	}
	if fileLogger.config.MaxAge != 7 {
		t.Errorf("Expected default MaxAge=7, got %d", fileLogger.config.MaxAge)
	}
}

// Simple mock logger for testing
type mockLogger struct {
	level Level
}

func (m *mockLogger) Debug(msg string, fields ...Field) {}
func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Warn(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Fatal(msg string, fields ...Field) {}
func (m *mockLogger) Panic(msg string, fields ...Field) {}
func (m *mockLogger) WithFields(fields ...Field) Logger {
	return m
}
func (m *mockLogger) WithContext(ctx context.Context) Logger {
	return m
}
