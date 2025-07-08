package db

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

func TestParseConfigDefaults(t *testing.T) {
	conf := viper.New()
	// No database config set, should use defaults

	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse config with defaults: %v", err)
	}

	if config.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected default port '3306', got '%s'", config.Port)
	}
	if config.User != "root" {
		t.Errorf("Expected default user 'root', got '%s'", config.User)
	}
	if config.Name != "scaffold" {
		t.Errorf("Expected default name 'scaffold', got '%s'", config.Name)
	}
	if config.MaxOpenConns != 25 {
		t.Errorf("Expected default max_open_conns 25, got %d", config.MaxOpenConns)
	}
	if config.MaxIdleConns != 5 {
		t.Errorf("Expected default max_idle_conns 5, got %d", config.MaxIdleConns)
	}
	if config.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("Expected default conn_max_lifetime 5m, got %v", config.ConnMaxLifetime)
	}
	if config.RetryAttempts != 5 {
		t.Errorf("Expected default retry_attempts 5, got %d", config.RetryAttempts)
	}
	if config.RetryDelay != 2*time.Second {
		t.Errorf("Expected default retry_delay 2s, got %v", config.RetryDelay)
	}
}

func TestBuildDSNWithEmptyPassword(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		Name:     "user",
	}

	expectedDSN := "root:@tcp(localhost:3306)/user?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	actualDSN := buildDSN(config)

	if actualDSN != expectedDSN {
		t.Errorf("Expected DSN '%s', got '%s'", expectedDSN, actualDSN)
	}
}

func TestParseConfigWithStructuredConfig(t *testing.T) {
	conf := viper.New()

	// Set up a structured database config using db.mysql (like in YAML files)
	dbConfig := map[string]interface{}{
		"host":               "127.0.0.1",
		"port":               "3306",
		"user":               "scaffold",
		"password":           "my_secure_password_123",
		"database":           "user",
		"max_open_conns":     30,
		"max_idle_conns":     8,
		"conn_max_lifetime":  "10m",
		"conn_max_idle_time": "3m",
		"retry_attempts":     3,
		"retry_delay":        "1s",
	}

	conf.Set("db.mysql", dbConfig)

	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse structured config: %v", err)
	}

	if config.Host != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected port '3306', got '%s'", config.Port)
	}
	if config.MaxOpenConns != 30 {
		t.Errorf("Expected max_open_conns 30, got %d", config.MaxOpenConns)
	}
	if config.RetryAttempts != 3 {
		t.Errorf("Expected retry_attempts 3, got %d", config.RetryAttempts)
	}
}

func TestParseConfigLegacySupport(t *testing.T) {
	conf := viper.New()

	// Test legacy "database" configuration still works (for backward compatibility)
	dbConfig := map[string]interface{}{
		"host":           "127.0.0.1",
		"port":           "3306",
		"user":           "scaffold",
		"password":       "legacy_password",
		"name":           "user",
		"max_open_conns": 40,
		"max_idle_conns": 12,
	}

	conf.Set("database", dbConfig)

	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse legacy config: %v", err)
	}

	if config.Host != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected port '3306', got '%s'", config.Port)
	}
	if config.MaxOpenConns != 40 {
		t.Errorf("Expected max_open_conns 40, got %d", config.MaxOpenConns)
	}
}

func createTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}

func TestNewConnectionInvalidConfig(t *testing.T) {
	conf := viper.New()
	logger := createTestLogger()

	// Test with invalid configuration that would cause connection to fail
	conf.Set("db.mysql.host", "nonexistent-host")
	conf.Set("db.mysql.port", "9999")
	conf.Set("db.mysql.retry_attempts", 1) // Reduce retry attempts for faster test
	conf.Set("db.mysql.retry_delay", "100ms")

	_, err := NewConnection(conf, logger)
	if err == nil {
		t.Error("Expected connection to fail with invalid host, but it succeeded")
	}
}
