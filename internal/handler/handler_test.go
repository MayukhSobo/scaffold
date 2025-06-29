package handler

import (
	"bytes"
	"scaffold/pkg/log"
	"testing"
)

func TestNewHandler(t *testing.T) {
	// Create a test logger
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	handler := NewHandler(logger)

	if handler == nil {
		t.Error("NewHandler() returned nil")
	}

	if handler.logger == nil {
		t.Error("Handler logger is nil")
	}
}

func TestGetLogger(t *testing.T) {
	// Create a test logger
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	handler := NewHandler(logger)
	retrievedLogger := handler.GetLogger()

	if retrievedLogger == nil {
		t.Error("GetLogger() returned nil")
	}

	// Test that the logger works
	retrievedLogger.Info("Test message")

	if buf.Len() == 0 {
		t.Error("Logger did not write output")
	}
}

func TestHandlerInterface(t *testing.T) {
	// Test that Handler satisfies its interface contract
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	handler := NewHandler(logger)

	// Ensure GetLogger method exists and returns correct type
	returnedLogger := handler.GetLogger()
	if returnedLogger != logger {
		t.Error("GetLogger() did not return the same logger instance")
	}
}
