package db

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

// Test with actual local config file loading
func TestWithActualLocalConfigFile(t *testing.T) {
	// Check if local.yml exists
	configPath := "../../configs/local.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Skipping test: local.yml config file not found")
	}

	// Load actual config file using Viper (same as the application)
	conf := viper.New()
	conf.SetConfigFile(configPath)
	err := conf.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config file %s: %v", configPath, err)
	}

	// Test that our database package can parse the real config
	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse real config file: %v", err)
	}

	// Verify that we get sensible values from the real config
	if config.Host == "" {
		t.Error("Host should not be empty from real config")
	}
	if config.Port == "" {
		t.Error("Port should not be empty from real config")
	}
	if config.User == "" {
		t.Error("User should not be empty from real config")
	}
	if config.Name == "" {
		t.Error("Database name should not be empty from real config")
	}

	// Test DSN generation with real config values
	dsn := buildDSN(config)
	if dsn == "" {
		t.Error("DSN should not be empty")
	}

	t.Logf("Successfully parsed real config - Host: %s, Port: %s, User: %s, DB: %s",
		config.Host, config.Port, config.User, config.Name)
	t.Logf("Generated DSN: %s", dsn)
}

func TestParseConfigLocalEnvironment(t *testing.T) {
	conf := viper.New()

	// Simulate local.yml configuration
	conf.Set("db.mysql.host", "127.0.0.1")
	conf.Set("db.mysql.port", "3306")
	conf.Set("db.mysql.user", "scaffold")
	conf.Set("db.mysql.password", "my_secure_password_123")
	conf.Set("db.mysql.database", "user")

	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse local config: %v", err)
	}

	if config.Host != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected port '3306', got '%s'", config.Port)
	}
	if config.User != "scaffold" {
		t.Errorf("Expected user 'scaffold', got '%s'", config.User)
	}
	if config.Password != "my_secure_password_123" {
		t.Errorf("Expected password 'my_secure_password_123', got '%s'", config.Password)
	}
	if config.Name != "user" {
		t.Errorf("Expected database name 'user', got '%s'", config.Name)
	}
}

func TestBuildDSNLocalEnvironment(t *testing.T) {
	// Test DSN building with local.yml values
	config := &Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "scaffold",
		Password: "my_secure_password_123",
		Name:     "user",
	}

	expectedDSN := "scaffold:my_secure_password_123@tcp(127.0.0.1:3306)/user?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	actualDSN := buildDSN(config)

	if actualDSN != expectedDSN {
		t.Errorf("Expected DSN '%s', got '%s'", expectedDSN, actualDSN)
	}
}

func TestNewConnectionLocalInvalidConfig(t *testing.T) {
	conf := viper.New()
	logger := createLocalTestLogger()

	// Test with invalid local configuration that would cause connection to fail
	conf.Set("db.mysql.host", "nonexistent-local-host")
	conf.Set("db.mysql.port", "9999")
	conf.Set("db.mysql.retry_attempts", 1) // Reduce retry attempts for faster test
	conf.Set("db.mysql.retry_delay", "100ms")

	_, err := NewConnection(conf, logger)
	if err == nil {
		t.Error("Expected connection to fail with invalid local host, but it succeeded")
	}
}

func createLocalTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}
