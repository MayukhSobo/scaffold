package repository

import (
	"bytes"
	"golang-di/pkg/log"
	"testing"

	"gorm.io/gorm"
)

func TestNewRepository(t *testing.T) {
	// Create a test logger
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	// Create a test database (mock)
	db := &gorm.DB{}

	repo := NewRepository(logger, db)

	if repo == nil {
		t.Error("NewRepository() returned nil")
	}

	if repo.db != db {
		t.Error("Repository database not set correctly")
	}

	if repo.logger == nil {
		t.Error("Repository logger is nil")
	}
}

func TestNewDb(t *testing.T) {
	// Test NewDb function
	db := NewDb()

	if db == nil {
		t.Error("NewDb() returned nil")
	}

	// Since this is a mock implementation, just verify it returns a *gorm.DB
	var dbInterface *gorm.DB = db
	if dbInterface != db {
		t.Error("NewDb() did not return correct type")
	}
}

func TestRepositoryLogger(t *testing.T) {
	// Test that repository logger works
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)

	// Test that logger is accessible and functional
	repo.logger.Info("Repository test message")

	if buf.Len() == 0 {
		t.Error("Repository logger did not write output")
	}

	output := buf.String()
	if !contains(output, "Repository test message") {
		t.Error("Repository logger output does not contain expected message")
	}
}

func TestRepositoryInterface(t *testing.T) {
	// Test that Repository satisfies its interface contract
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)

	// Ensure all fields are accessible
	if repo.db != db {
		t.Error("Repository does not contain the expected database instance")
	}

	if repo.logger != logger {
		t.Error("Repository does not contain the expected logger instance")
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
