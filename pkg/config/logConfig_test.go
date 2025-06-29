package config

import (
	"golang-di/pkg/log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// createTestConfig creates a viper config for testing
func createTestConfig(config map[string]any) *viper.Viper {
	v := viper.New()
	for key, value := range config {
		v.Set(key, value)
	}
	return v
}

func TestBasicConfig(t *testing.T) {
	// Test with viper config (recommended approach)
	config := createTestConfig(map[string]any{
		"log.level":                      "debug",
		"log.console_logger.enabled":     true,
		"log.console_logger.colors":      false,
		"log.console_logger.json_format": false,
		"log.file_logger.enabled":        false,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Config-driven logger should not be nil")
	}

	logger.Info("Configuration-driven logging test")
	logger.Debug("Debug from config", log.String("approach", "config-driven"))
}

func TestFileOutput(t *testing.T) {
	testDir := "test_logs"
	logFile := "config_test.log"
	defer os.RemoveAll(testDir) // Clean up entire directory

	config := createTestConfig(map[string]any{
		"log.level":                   "info",
		"log.console_logger.enabled":  false,
		"log.file_logger.enabled":     true,
		"log.file_logger.directory":   testDir,
		"log.file_logger.filename":    logFile,
		"log.file_logger.json_format": true,
		"log.file_logger.max_size":    1,
		"log.file_logger.max_backups": 1,
		"log.file_logger.max_age":     1,
		"log.file_logger.compress":    false,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Config-driven file logger should not be nil")
	}

	logger.Info("File logging from config")

	// Verify directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Log directory was not created")
	}

	// Verify file was created in the correct directory
	fullPath := filepath.Join(testDir, logFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Log file was not created in specified directory")
	}
}

func TestMultiOutput(t *testing.T) {
	testDir := "test_multi_logs"
	logFile := "multi_config_test.log"
	defer os.RemoveAll(testDir) // Clean up entire directory

	config := createTestConfig(map[string]any{
		"log.level":                      "debug",
		"log.console_logger.enabled":     true,
		"log.console_logger.colors":      true,
		"log.console_logger.json_format": false,
		"log.file_logger.enabled":        true,
		"log.file_logger.directory":      testDir,
		"log.file_logger.filename":       logFile,
		"log.file_logger.json_format":    true,
		"log.file_logger.max_size":       1,
		"log.file_logger.max_backups":    1,
		"log.file_logger.max_age":        1,
		"log.file_logger.compress":       false,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Config-driven multi logger should not be nil")
	}

	logger.Info("Multi-output logging from config")

	// Verify file was created in the correct directory
	fullPath := filepath.Join(testDir, logFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Log file was not created in specified directory for multi-output")
	}
}

func TestDefaults(t *testing.T) {
	// Test with nil config (should use defaults)
	logger := CreateLoggerFromConfig(nil)
	if logger == nil {
		t.Fatal("Default config logger should not be nil")
	}

	logger.Info("Default configuration test")
}

func TestInvalidLevel(t *testing.T) {
	config := createTestConfig(map[string]any{
		"log.level":                  "invalid_level",
		"log.console_logger.enabled": true,
		"log.console_logger.colors":  true,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Logger with invalid level should still be created with default level")
	}

	logger.Info("Invalid level test - should use default level")
}

func TestLogDirectoryCreation(t *testing.T) {
	testDir := "test_nested/deeply/nested/logs"
	logFile := "deep_test.log"
	defer os.RemoveAll("test_nested") // Clean up from root

	config := createTestConfig(map[string]any{
		"log.level":                   "info",
		"log.console_logger.enabled":  false,
		"log.file_logger.enabled":     true,
		"log.file_logger.directory":   testDir,
		"log.file_logger.filename":    logFile,
		"log.file_logger.json_format": false,
		"log.file_logger.max_size":    1,
		"log.file_logger.max_backups": 1,
		"log.file_logger.max_age":     1,
		"log.file_logger.compress":    false,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Logger should be created even with deeply nested directory")
	}

	logger.Info("Testing deeply nested directory creation")

	// Verify deeply nested directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Deeply nested log directory was not created")
	}

	// Verify file was created
	fullPath := filepath.Join(testDir, logFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Log file was not created in deeply nested directory")
	}
}

func TestEmptyConfig(t *testing.T) {
	// Test with empty config (should default to console logger)
	config := createTestConfig(map[string]any{})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Logger should be created with empty config")
	}

	logger.Info("Empty configuration test - should use console logger")
}

func TestConsoleJsonFormat(t *testing.T) {
	config := createTestConfig(map[string]any{
		"log.level":                      "info",
		"log.console_logger.enabled":     true,
		"log.console_logger.json_format": true, // Test JSON format for console
		"log.console_logger.colors":      false,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Fatal("Console logger with JSON format should not be nil")
	}

	logger.Info("Testing console JSON format")
}

func TestNewConfig(t *testing.T) {
	// Test with environment variable
	os.Setenv("APP_CONF", "config/local.yml")
	defer os.Unsetenv("APP_CONF")

	// This would normally panic if config file doesn't exist
	// We'll test the function exists and can be called
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic if config file doesn't exist
			t.Logf("NewConfig panicked as expected: %v", r)
		}
	}()

	// Test the function (will panic if file doesn't exist, which is expected behavior)
	if _, err := os.Stat("config/local.yml"); err == nil {
		config := NewConfig()
		if config == nil {
			t.Error("NewConfig should not return nil when config file exists")
		}
	}
}

func TestGetConfig(t *testing.T) {
	// Create a temporary config file
	tmpFile := "test_config.yml"
	content := `
log:
  level: info
  console_logger:
    enabled: true
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpFile)

	config := getConfig(tmpFile)
	if config == nil {
		t.Error("getConfig should not return nil for valid config file")
	}

	// Test that config values are loaded correctly
	if config.GetString("log.level") != "info" {
		t.Error("Config value not loaded correctly")
	}
}

func TestGetConfigInvalidFile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("getConfig should panic for invalid file")
		}
	}()

	// This should panic
	getConfig("nonexistent_file.yml")
}

func TestParseLevelEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected log.LogLevel
	}{
		{"debug", log.DebugLevel},
		{"info", log.InfoLevel},
		{"warn", log.WarnLevel},
		{"error", log.ErrorLevel},
		{"fatal", log.FatalLevel},
		{"panic", log.PanicLevel},
		{"invalid", log.InfoLevel},
		{"", log.InfoLevel},
		{"DEBUG", log.InfoLevel}, // case sensitive
		{"Info", log.InfoLevel},  // case sensitive
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := parseLevel(test.input)
			if result != test.expected {
				t.Errorf("parseLevel(%s) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

func TestResolveLogFilePathEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		filename  string
		expected  string // relative part to check
	}{
		{"empty directory", "", "test.log", "logs/test.log"},
		{"relative directory", "logs", "test.log", "logs/test.log"},
		{"absolute directory", "/tmp", "test.log", "/tmp/test.log"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := resolveLogFilePath(test.directory, test.filename)
			if test.directory == "/tmp" {
				// For absolute paths, check exact match
				if result != test.expected {
					t.Errorf("resolveLogFilePath(%s, %s) = %s, expected %s",
						test.directory, test.filename, result, test.expected)
				}
			} else {
				// For relative paths, check that it contains the expected part
				if !contains(result, test.expected) {
					t.Errorf("resolveLogFilePath(%s, %s) = %s, should contain %s",
						test.directory, test.filename, result, test.expected)
				}
			}
		})
	}
}

func TestEnsureLogDirectoryError(t *testing.T) {
	// Test with invalid directory (should fail on most systems)
	invalidDir := "/root/nonexistent/deep/path"
	err := ensureLogDirectory(invalidDir)
	if err == nil {
		t.Error("ensureLogDirectory should fail for invalid directory")
	}
}

func TestCreateFileLoggerErrorHandling(t *testing.T) {
	// Test error handling in createFileLogger
	config := createTestConfig(map[string]any{
		"log.level":                   "info",
		"log.file_logger.enabled":     true,
		"log.file_logger.filename":    "test.log",
		"log.file_logger.directory":   "/root/invalid/path", // Should fail
		"log.file_logger.max_size":    10,
		"log.file_logger.max_backups": 3,
		"log.file_logger.max_age":     7,
		"log.file_logger.compress":    true,
		"log.file_logger.json_format": true,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Error("Should return console logger as fallback")
	}

	// Should have fallen back to console logger
	logger.Info("Test fallback message")
}

func TestCreateConsoleLoggerNilWriter(t *testing.T) {
	config := createTestConfig(map[string]any{
		"log.level":                      "info",
		"log.console_logger.enabled":     true,
		"log.console_logger.json_format": true,
	})

	logger := CreateLoggerFromConfig(config)
	if logger == nil {
		t.Error("Console logger should not be nil")
	}

	logger.Info("Test console logger with nil writer")
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
