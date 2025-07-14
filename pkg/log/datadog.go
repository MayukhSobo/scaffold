package log

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// DatadogLoggerConfig contains configuration for Datadog logging.
type DatadogLoggerConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Service     string `mapstructure:"service"`
	Environment string `mapstructure:"environment"`
	Source      string `mapstructure:"source"`
	Tags        string `mapstructure:"tags"`
	Timeout     int    `mapstructure:"timeout"`     // timeout in seconds for connection
	JsonFormat  bool   `mapstructure:"json_format"` // whether to use JSON format
}

// DatadogLogger implements Logger interface for Datadog output via TCP.
type DatadogLogger struct {
	config      *DatadogLoggerConfig
	level       Level
	contextData map[string]any
	conn        net.Conn
	connMutex   sync.RWMutex
	address     string
}

// DatadogLogEntry represents a log entry in JSON format for Datadog.
type DatadogLogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	Environment string                 `json:"environment"`
	Source      string                 `json:"source"`
	Tags        string                 `json:"tags,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// preparedLogData holds all the prepared log data for formatting.
type preparedLogData struct {
	Timestamp   string
	Level       string
	Message     string
	Service     string
	Environment string
	Source      string
	Tags        string
	Fields      map[string]interface{}
}

func init() {
	RegisterFactory("datadog", NewDatadogLoggerFromConfig)
}

// NewDatadogLoggerFromConfig creates a new Datadog logger from a Viper configuration.
func NewDatadogLoggerFromConfig(level Level, v *viper.Viper) (Logger, error) {
	var config DatadogLoggerConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal datadog logger config: %w", err)
	}

	// Set defaults
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == 0 {
		config.Port = 10518
	}
	if config.Source == "" {
		config.Source = "go"
	}
	if config.Service == "" {
		config.Service = "scaffold"
	}
	if config.Environment == "" {
		config.Environment = "development"
	}
	if config.Timeout == 0 {
		config.Timeout = 5 // 5 seconds default timeout
	}

	return NewDatadogLogger(level, &config), nil
}

// NewDatadogLogger creates a new Datadog logger.
func NewDatadogLogger(level Level, config *DatadogLoggerConfig) Logger {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	return &DatadogLogger{
		config:      config,
		level:       level,
		contextData: make(map[string]any),
		address:     address,
	}
}

// ensureConnection ensures we have a valid TCP connection to the Datadog agent.
func (d *DatadogLogger) ensureConnection() error {
	d.connMutex.RLock()
	if d.conn != nil {
		d.connMutex.RUnlock()
		return nil
	}
	d.connMutex.RUnlock()

	d.connMutex.Lock()
	defer d.connMutex.Unlock()

	// Double-check after acquiring write lock
	if d.conn != nil {
		return nil
	}

	// Create connection with timeout
	conn, err := net.DialTimeout("tcp", d.address, time.Duration(d.config.Timeout)*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to Datadog agent at %s: %w", d.address, err)
	}

	d.conn = conn
	return nil
}

// sendLogEntry sends a log entry to Datadog agent over TCP.
func (d *DatadogLogger) sendLogEntry(level, message string, fields []Field) {
	// Build structured log line
	logLine := d.buildLogLine(level, message, fields)

	// Send asynchronously to avoid blocking
	go func() {
		if err := d.ensureConnection(); err != nil {
			// If we can't connect, silently fail to avoid logging loops
			return
		}

		d.connMutex.RLock()
		conn := d.conn
		d.connMutex.RUnlock()

		if conn != nil {
			// Set write deadline to prevent hanging
			conn.SetWriteDeadline(time.Now().Add(time.Duration(d.config.Timeout) * time.Second))

			_, err := conn.Write([]byte(logLine + "\n"))
			if err != nil {
				// Connection failed, close it and next log will try to reconnect
				d.connMutex.Lock()
				if d.conn != nil {
					d.conn.Close()
					d.conn = nil
				}
				d.connMutex.Unlock()
			}
		}
	}()
}

// buildLogLine creates a structured log line in either text or JSON format for Datadog.
func (d *DatadogLogger) buildLogLine(level, message string, fields []Field) string {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	if d.config.JsonFormat {
		return d.jsonify(timestamp, level, message, fields)
	}
	return d.fancy(timestamp, level, message, fields)
}

// processLogs collects all log metadata and fields into a structured format.
func (d *DatadogLogger) processLogs(timestamp, level, message string, fields []Field) *preparedLogData {
	// Collect all fields (context + provided)
	allFields := make(map[string]interface{})

	// Add context data
	for k, v := range d.contextData {
		allFields[k] = v
	}

	// Add provided fields
	for _, field := range fields {
		allFields[field.Key] = field.Value
	}

	return &preparedLogData{
		Timestamp:   timestamp,
		Level:       level,
		Message:     message,
		Service:     d.config.Service,
		Environment: d.config.Environment,
		Source:      d.config.Source,
		Tags:        d.config.Tags,
		Fields:      allFields,
	}
}

// jsonify creates a JSON-formatted log line.
func (d *DatadogLogger) jsonify(timestamp, level, message string, fields []Field) string {
	data := d.processLogs(timestamp, level, message, fields)

	entry := DatadogLogEntry{
		Timestamp:   data.Timestamp,
		Level:       data.Level,
		Message:     data.Message,
		Service:     data.Service,
		Environment: data.Environment,
		Source:      data.Source,
		Tags:        data.Tags,
		Fields:      data.Fields,
	}

	// Remove empty fields map if no fields
	if len(entry.Fields) == 0 {
		entry.Fields = nil
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		// If we can't marshal, fall back to text format
		return d.fancy(timestamp, level, message, fields)
	}

	return string(jsonData)
}

// fancy creates a text-formatted fancy log line.
func (d *DatadogLogger) fancy(timestamp, level, message string, fields []Field) string {
	data := d.processLogs(timestamp, level, message, fields)

	// Start with basic log format
	logLine := fmt.Sprintf("%s %s service=%s env=%s source=%s",
		data.Timestamp,
		data.Level,
		data.Service,
		data.Environment,
		data.Source)

	// Add tags if configured
	if data.Tags != "" {
		logLine += fmt.Sprintf(" tags=%s", data.Tags)
	}

	// Add all fields
	for k, v := range data.Fields {
		logLine += fmt.Sprintf(" %s=%v", k, v)
	}

	// Add the actual message
	logLine += fmt.Sprintf(" msg=\"%s\"", data.Message)

	return logLine
}

// Debug logs a debug message.
func (d *DatadogLogger) Debug(msg string, fields ...Field) {
	d.sendLogEntry("DEBUG", msg, fields)
}

// Info logs an info message.
func (d *DatadogLogger) Info(msg string, fields ...Field) {
	d.sendLogEntry("INFO", msg, fields)
}

// Warn logs a warning message.
func (d *DatadogLogger) Warn(msg string, fields ...Field) {
	d.sendLogEntry("WARN", msg, fields)
}

// Error logs an error message.
func (d *DatadogLogger) Error(msg string, fields ...Field) {
	d.sendLogEntry("ERROR", msg, fields)
}

// Fatal logs a fatal message.
func (d *DatadogLogger) Fatal(msg string, fields ...Field) {
	d.sendLogEntry("FATAL", msg, fields)
}

// Panic logs a panic message.
func (d *DatadogLogger) Panic(msg string, fields ...Field) {
	d.sendLogEntry("PANIC", msg, fields)
}

// Formatted logging methods
func (d *DatadogLogger) Debugf(format string, args ...interface{}) {
	d.Debug(fmt.Sprintf(format, args...))
}

func (d *DatadogLogger) Infof(format string, args ...interface{}) {
	d.Info(fmt.Sprintf(format, args...))
}

func (d *DatadogLogger) Warnf(format string, args ...interface{}) {
	d.Warn(fmt.Sprintf(format, args...))
}

func (d *DatadogLogger) Errorf(format string, args ...interface{}) {
	d.Error(fmt.Sprintf(format, args...))
}

func (d *DatadogLogger) Fatalf(format string, args ...interface{}) {
	d.Fatal(fmt.Sprintf(format, args...))
}

func (d *DatadogLogger) Panicf(format string, args ...interface{}) {
	d.Panic(fmt.Sprintf(format, args...))
}

// WithFields creates a new logger with additional context fields.
func (d *DatadogLogger) WithFields(fields ...Field) Logger {
	newContextData := make(map[string]any)

	// Copy existing context data
	for k, v := range d.contextData {
		newContextData[k] = v
	}

	// Add new fields
	for _, field := range fields {
		newContextData[field.Key] = field.Value
	}

	return &DatadogLogger{
		config:      d.config,
		level:       d.level,
		contextData: newContextData,
		conn:        d.conn, // Share connection
		address:     d.address,
	}
}

// WithContext creates a new logger with context.
func (d *DatadogLogger) WithContext(ctx context.Context) Logger {
	return &DatadogLogger{
		config:      d.config,
		level:       d.level,
		contextData: d.contextData,
		conn:        d.conn, // Share connection
		address:     d.address,
	}
}

// Close closes the TCP connection to the Datadog agent.
func (d *DatadogLogger) Close() error {
	d.connMutex.Lock()
	defer d.connMutex.Unlock()

	if d.conn != nil {
		err := d.conn.Close()
		d.conn = nil
		return err
	}
	return nil
}
