package service

import (
	"bytes"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"testing"
)

func TestNewService(t *testing.T) {
	// Create a test logger
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	if service == nil {
		t.Error("NewService() returned nil")
	}

	if service.logger == nil {
		t.Error("Service logger is nil")
	}
}

func TestServiceLogger(t *testing.T) {
	// Create a test logger
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	// Test that the logger works
	service.logger.Info("Service test message")

	if buf.Len() == 0 {
		t.Error("Service logger did not write output")
	}

	output := buf.String()
	if !contains(output, "Service test message") {
		t.Error("Service logger output does not contain expected message")
	}
}

func TestServiceInterface(t *testing.T) {
	// Test that Service satisfies its interface contract
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	// Ensure logger field is accessible
	if service.logger != logger {
		t.Error("Service does not contain the expected logger instance")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
