package db

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

// Test with actual Docker config file loading
func TestWithActualDockerConfigFile(t *testing.T) {
	// Check if docker.yml exists
	configPath := "../../configs/docker.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Skipping test: docker.yml config file not found")
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

func TestWithActualDockerConfigFileConnection(t *testing.T) {
	// Check if docker.yml exists
	configPath := "../../configs/docker.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Skipping test: docker.yml config file not found")
	}

	// Load actual config file using Viper (same as the application)
	conf := viper.New()
	conf.SetConfigFile(configPath)
	err := conf.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config file %s: %v", configPath, err)
	}

	// Set reduced retry for faster test
	conf.Set("db.mysql.retry_attempts", 1)
	conf.Set("db.mysql.retry_delay", "100ms")

	logger := createDockerTestLogger()

	// Try to connect using the real Docker config - should fail unless Docker is running
	_, err = NewConnection(conf, logger)
	if err == nil {
		t.Error("Expected connection to fail with docker.yml config (mysql hostname), but it succeeded")
	}

	t.Logf("Connection failed as expected: %v", err)
}

func TestParseConfigDockerEnvironment(t *testing.T) {
	conf := viper.New()

	// Simulate docker.yml configuration
	conf.Set("db.mysql.host", "mysql")
	conf.Set("db.mysql.port", "3306")
	conf.Set("db.mysql.user", "scaffold")
	conf.Set("db.mysql.password", "bXlfc2VjdXJlX3Bhc3N3b3JkXzEyMw==") // base64 encoded
	conf.Set("db.mysql.database", "user")

	config, err := parseConfig(conf)
	if err != nil {
		t.Fatalf("Failed to parse docker config: %v", err)
	}

	if config.Host != "mysql" {
		t.Errorf("Expected host 'mysql', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected port '3306', got '%s'", config.Port)
	}
	if config.User != "scaffold" {
		t.Errorf("Expected user 'scaffold', got '%s'", config.User)
	}
	if config.Password != "bXlfc2VjdXJlX3Bhc3N3b3JkXzEyMw==" {
		t.Errorf("Expected encoded password, got '%s'", config.Password)
	}
	if config.Name != "user" {
		t.Errorf("Expected database name 'user', got '%s'", config.Name)
	}
}

func TestBuildDSNDockerEnvironment(t *testing.T) {
	// Test DSN building with docker.yml values
	config := &Config{
		Host:     "mysql",
		Port:     "3306",
		User:     "scaffold",
		Password: "bXlfc2VjdXJlX3Bhc3N3b3JkXzEyMw==",
		Name:     "user",
	}

	expectedDSN := "scaffold:bXlfc2VjdXJlX3Bhc3N3b3JkXzEyMw==@tcp(mysql:3306)/user?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	actualDSN := buildDSN(config)

	if actualDSN != expectedDSN {
		t.Errorf("Expected DSN '%s', got '%s'", expectedDSN, actualDSN)
	}
}

func TestNewConnectionDockerInvalidConfig(t *testing.T) {
	conf := viper.New()
	logger := createDockerTestLogger()

	// Test with invalid Docker configuration that would cause connection to fail
	conf.Set("db.mysql.host", "nonexistent-docker-host")
	conf.Set("db.mysql.port", "9999")
	conf.Set("db.mysql.retry_attempts", 1) // Reduce retry attempts for faster test
	conf.Set("db.mysql.retry_delay", "100ms")

	_, err := NewConnection(conf, logger)
	if err == nil {
		t.Error("Expected connection to fail with invalid Docker host, but it succeeded")
	}
}

func TestNewConnectionDockerRealConfig(t *testing.T) {
	conf := viper.New()
	logger := createDockerTestLogger()

	// Test with actual Docker config (mysql hostname) - should fail unless Docker is running
	conf.Set("db.mysql.host", "mysql")
	conf.Set("db.mysql.port", "3306")
	conf.Set("db.mysql.user", "scaffold")
	conf.Set("db.mysql.password", "bXlfc2VjdXJlX3Bhc3N3b3JkXzEyMw==")
	conf.Set("db.mysql.database", "user")
	conf.Set("db.mysql.retry_attempts", 1) // Reduce retry attempts for faster test
	conf.Set("db.mysql.retry_delay", "100ms")

	_, err := NewConnection(conf, logger)
	if err == nil {
		t.Error("Expected connection to fail with 'mysql' hostname (unless Docker is running), but it succeeded")
	}
}

func createDockerTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}
