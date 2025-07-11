package log

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestDatadogLoggerCreation(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:        "127.0.0.1",
		Port:        10518,
		Service:     "test-service",
		Environment: "test",
		Source:      "go",
		Tags:        "env:test,service:test",
		Timeout:     5,
		JsonFormat:  true,
	}

	logger := NewDatadogLogger(InfoLevel, config)
	if logger == nil {
		t.Fatal("Datadog logger should not be nil")
	}

	// Cast to get access to internal config
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	if datadogLogger.config.Service != "test-service" {
		t.Errorf("Expected service='test-service', got '%s'", datadogLogger.config.Service)
	}

	if datadogLogger.config.Port != 10518 {
		t.Errorf("Expected port=10518, got %d", datadogLogger.config.Port)
	}

	if datadogLogger.address != "127.0.0.1:10518" {
		t.Errorf("Expected address='127.0.0.1:10518', got '%s'", datadogLogger.address)
	}

	if !datadogLogger.config.JsonFormat {
		t.Errorf("Expected json_format=true, got %v", datadogLogger.config.JsonFormat)
	}
}

func TestDatadogLoggerFromConfig(t *testing.T) {
	v := viper.New()
	v.Set("host", "127.0.0.1")
	v.Set("port", 10518)
	v.Set("service", "test-service")
	v.Set("environment", "test")
	v.Set("source", "go")
	v.Set("tags", "env:test")
	v.Set("timeout", 10)
	v.Set("json_format", true)

	logger, err := NewDatadogLoggerFromConfig(InfoLevel, v)
	if err != nil {
		t.Fatalf("Failed to create Datadog logger from config: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	if datadogLogger.config.Host != "127.0.0.1" {
		t.Errorf("Expected host='127.0.0.1', got '%s'", datadogLogger.config.Host)
	}

	if datadogLogger.config.Timeout != 10 {
		t.Errorf("Expected timeout=10, got %d", datadogLogger.config.Timeout)
	}

	if !datadogLogger.config.JsonFormat {
		t.Errorf("Expected json_format=true, got %v", datadogLogger.config.JsonFormat)
	}
}

func TestDatadogLoggerDefaults(t *testing.T) {
	v := viper.New()
	// Only set minimal required config

	logger, err := NewDatadogLoggerFromConfig(InfoLevel, v)
	if err != nil {
		t.Fatalf("Failed to create Datadog logger with defaults: %v", err)
	}

	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	// Test defaults
	if datadogLogger.config.Host != "127.0.0.1" {
		t.Errorf("Expected default host='127.0.0.1', got '%s'", datadogLogger.config.Host)
	}

	if datadogLogger.config.Port != 10518 {
		t.Errorf("Expected default port=10518, got %d", datadogLogger.config.Port)
	}

	if datadogLogger.config.Service != "scaffold" {
		t.Errorf("Expected default service='scaffold', got '%s'", datadogLogger.config.Service)
	}

	if datadogLogger.config.Source != "go" {
		t.Errorf("Expected default source='go', got '%s'", datadogLogger.config.Source)
	}

	if datadogLogger.config.Timeout != 5 {
		t.Errorf("Expected default timeout=5, got %d", datadogLogger.config.Timeout)
	}

	if datadogLogger.config.JsonFormat {
		t.Errorf("Expected default json_format=false, got %v", datadogLogger.config.JsonFormat)
	}

	if datadogLogger.address != "127.0.0.1:10518" {
		t.Errorf("Expected address='127.0.0.1:10518', got '%s'", datadogLogger.address)
	}
}

func TestDatadogLoggerInterface(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:       "127.0.0.1",
		Port:       10518,
		Service:    "test-service",
		Timeout:    1, // Short timeout for testing
		JsonFormat: true,
	}

	logger := NewDatadogLogger(InfoLevel, config)

	// Test all interface methods (they should not panic)
	logger.Info("Test info message")
	logger.Debug("Test debug message", String("key", "value"))
	logger.Warn("Test warning message")
	logger.Error("Test error message", Error(nil))

	// Test formatted methods
	logger.Infof("Test formatted info: %s", "test")
	logger.Debugf("Test formatted debug: %d", 42)
	logger.Warnf("Test formatted warning: %.2f", 3.14)
	logger.Errorf("Test formatted error: %v", "error")

	// Test WithFields
	contextLogger := logger.WithFields(
		String("request_id", "test-123"),
		String("user_id", "user-456"),
	)

	contextLogger.Info("Test with context fields")

	// Test WithContext
	ctxLogger := logger.WithContext(nil)
	ctxLogger.Info("Test with context")
}

func TestDatadogLoggerRegistration(t *testing.T) {
	// Test that the factory was registered
	factory, ok := loggerFactories["datadog"]
	if !ok {
		t.Fatal("Datadog logger factory not registered")
	}

	v := viper.New()
	v.Set("service", "test-service")
	v.Set("host", "127.0.0.1")
	v.Set("port", 10518)
	v.Set("json_format", true)

	logger, err := factory(InfoLevel, v)
	if err != nil {
		t.Fatalf("Factory failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Factory returned nil logger")
	}

	// Verify it's a DatadogLogger
	_, ok = logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Factory did not return a DatadogLogger")
	}
}

func TestDatadogLoggerConcurrency(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:       "127.0.0.1",
		Port:       10518,
		Service:    "test-service",
		Timeout:    1,
		JsonFormat: true,
	}

	logger := NewDatadogLogger(InfoLevel, config)

	// Test concurrent logging
	done := make(chan bool)
	numGoroutines := 10
	messagesPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info("Concurrent test message",
					String("goroutine_id", string(rune(id))),
					Int("message_id", j),
				)
				time.Sleep(time.Millisecond)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// If we reach here without panic, the test passed
	t.Log("Concurrent logging test completed successfully")
}

func TestDatadogLoggerBuildLogLineText(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:        "127.0.0.1",
		Port:        10518,
		Service:     "test-service",
		Environment: "test",
		Source:      "go",
		Tags:        "env:test,service:test",
		JsonFormat:  false, // Test text format
	}

	logger := NewDatadogLogger(InfoLevel, config)
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	fields := []Field{
		String("user_id", "123"),
		Int("count", 42),
	}

	logLine := datadogLogger.buildLogLine("INFO", "Test message", fields)

	// Check that the log line contains expected components
	if !contains(logLine, "INFO") {
		t.Error("Log line should contain log level")
	}

	if !contains(logLine, "service=test-service") {
		t.Error("Log line should contain service")
	}

	if !contains(logLine, "env=test") {
		t.Error("Log line should contain environment")
	}

	if !contains(logLine, "source=go") {
		t.Error("Log line should contain source")
	}

	if !contains(logLine, "tags=env:test,service:test") {
		t.Error("Log line should contain tags")
	}

	if !contains(logLine, "user_id=123") {
		t.Error("Log line should contain user_id field")
	}

	if !contains(logLine, "count=42") {
		t.Error("Log line should contain count field")
	}

	if !contains(logLine, `msg="Test message"`) {
		t.Error("Log line should contain the message")
	}
}

func TestDatadogLoggerBuildLogLineJSON(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:        "127.0.0.1",
		Port:        10518,
		Service:     "test-service",
		Environment: "test",
		Source:      "go",
		Tags:        "env:test,service:test",
		JsonFormat:  true, // Test JSON format
	}

	logger := NewDatadogLogger(InfoLevel, config)
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	fields := []Field{
		String("user_id", "123"),
		Int("count", 42),
	}

	logLine := datadogLogger.buildLogLine("INFO", "Test message", fields)

	// Parse the JSON to verify it's valid JSON
	var entry DatadogLogEntry
	err := json.Unmarshal([]byte(logLine), &entry)
	if err != nil {
		t.Fatalf("Failed to parse JSON log line: %v", err)
	}

	// Check JSON structure
	if entry.Level != "INFO" {
		t.Errorf("Expected level='INFO', got '%s'", entry.Level)
	}

	if entry.Message != "Test message" {
		t.Errorf("Expected message='Test message', got '%s'", entry.Message)
	}

	if entry.Service != "test-service" {
		t.Errorf("Expected service='test-service', got '%s'", entry.Service)
	}

	if entry.Environment != "test" {
		t.Errorf("Expected environment='test', got '%s'", entry.Environment)
	}

	if entry.Source != "go" {
		t.Errorf("Expected source='go', got '%s'", entry.Source)
	}

	if entry.Tags != "env:test,service:test" {
		t.Errorf("Expected tags='env:test,service:test', got '%s'", entry.Tags)
	}

	// Check fields
	if entry.Fields["user_id"] != "123" {
		t.Errorf("Expected user_id='123', got '%v'", entry.Fields["user_id"])
	}

	// JSON unmarshaling converts numbers to float64
	if entry.Fields["count"] != float64(42) {
		t.Errorf("Expected count=42, got '%v'", entry.Fields["count"])
	}

	// Verify timestamp is parseable
	_, err = time.Parse(time.RFC3339, entry.Timestamp)
	if err != nil {
		t.Errorf("Timestamp should be valid RFC3339: %v", err)
	}
}

func TestDatadogLoggerJSONFormatToggle(t *testing.T) {
	// Test that we can toggle between JSON and text formats
	config := &DatadogLoggerConfig{
		Host:        "127.0.0.1",
		Port:        10518,
		Service:     "test-service",
		Environment: "test",
		Source:      "go",
		Tags:        "env:test",
		JsonFormat:  false,
	}

	logger := NewDatadogLogger(InfoLevel, config)
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	fields := []Field{String("key", "value")}

	// Test text format
	textLogLine := datadogLogger.buildLogLine("INFO", "Test message", fields)
	if contains(textLogLine, "{") || contains(textLogLine, "}") {
		t.Error("Text format should not contain JSON braces")
	}

	// Switch to JSON format
	datadogLogger.config.JsonFormat = true
	jsonLogLine := datadogLogger.buildLogLine("INFO", "Test message", fields)

	// Should be valid JSON
	var entry DatadogLogEntry
	err := json.Unmarshal([]byte(jsonLogLine), &entry)
	if err != nil {
		t.Errorf("JSON format should produce valid JSON: %v", err)
	}
}

func TestDatadogLoggerClose(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:    "127.0.0.1",
		Port:    10518,
		Service: "test-service",
	}

	logger := NewDatadogLogger(InfoLevel, config)
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	// Close should not error even if no connection was established
	err := datadogLogger.Close()
	if err != nil {
		t.Errorf("Close should not error: %v", err)
	}
}

func TestDatadogLoggerProcessLogs(t *testing.T) {
	config := &DatadogLoggerConfig{
		Host:        "127.0.0.1",
		Port:        10518,
		Service:     "test-service",
		Environment: "test",
		Source:      "go",
		Tags:        "env:test,service:test",
	}

	logger := NewDatadogLogger(InfoLevel, config)
	datadogLogger, ok := logger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast to *DatadogLogger")
	}

	// Add some context data
	contextLogger := datadogLogger.WithFields(
		String("request_id", "req-123"),
		String("user_id", "user-456"),
	)

	contextDatadogLogger, ok := contextLogger.(*DatadogLogger)
	if !ok {
		t.Fatal("Could not cast context logger to *DatadogLogger")
	}

	// Prepare log data
	fields := []Field{
		String("action", "login"),
		Int("attempts", 3),
	}

	data := contextDatadogLogger.processLogs("2024-01-15T12:30:45Z", "INFO", "Test message", fields)

	// Verify the prepared data
	if data.Timestamp != "2024-01-15T12:30:45Z" {
		t.Errorf("Expected timestamp='2024-01-15T12:30:45Z', got '%s'", data.Timestamp)
	}

	if data.Level != "INFO" {
		t.Errorf("Expected level='INFO', got '%s'", data.Level)
	}

	if data.Message != "Test message" {
		t.Errorf("Expected message='Test message', got '%s'", data.Message)
	}

	if data.Service != "test-service" {
		t.Errorf("Expected service='test-service', got '%s'", data.Service)
	}

	if data.Environment != "test" {
		t.Errorf("Expected environment='test', got '%s'", data.Environment)
	}

	if data.Source != "go" {
		t.Errorf("Expected source='go', got '%s'", data.Source)
	}

	if data.Tags != "env:test,service:test" {
		t.Errorf("Expected tags='env:test,service:test', got '%s'", data.Tags)
	}

	// Verify fields (context + provided)
	if data.Fields["request_id"] != "req-123" {
		t.Errorf("Expected request_id='req-123', got '%v'", data.Fields["request_id"])
	}

	if data.Fields["user_id"] != "user-456" {
		t.Errorf("Expected user_id='user-456', got '%v'", data.Fields["user_id"])
	}

	if data.Fields["action"] != "login" {
		t.Errorf("Expected action='login', got '%v'", data.Fields["action"])
	}

	if data.Fields["attempts"] != 3 {
		t.Errorf("Expected attempts=3, got '%v'", data.Fields["attempts"])
	}

	// Verify all fields are present
	expectedFieldCount := 4 // request_id, user_id, action, attempts
	if len(data.Fields) != expectedFieldCount {
		t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(data.Fields))
	}
}
