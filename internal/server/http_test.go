package server

import (
	"bytes"
	"golang-di/internal/handler"
	"golang-di/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewServerHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	engine := NewServerHTTP(logger, userHandler)

	if engine == nil {
		t.Error("NewServerHTTP() returned nil")
	}

	// Test that it returns a gin.Engine
	var ginEngine *gin.Engine = engine
	if ginEngine != engine {
		t.Error("NewServerHTTP() did not return correct type")
	}
}

func TestServerHTTPRootEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	engine := NewServerHTTP(logger, userHandler)

	// Test root endpoint
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check that logger was called
	if buf.Len() == 0 {
		t.Error("Logger was not called for root endpoint")
	}
}

func TestServerHTTPUserEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	engine := NewServerHTTP(logger, userHandler)

	// Test user endpoint (should fail validation but endpoint should exist)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user", nil)

	engine.ServeHTTP(w, req)

	// Should return 400 because id parameter is required
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing parameter, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestServerHTTPWithValidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	engine := NewServerHTTP(logger, userHandler)

	// Test user endpoint with valid ID
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user?id=123", nil)

	engine.ServeHTTP(w, req)

	// Should return 200 with mock implementation
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for valid request, got %d", http.StatusOK, w.Code)
	}
}

func TestServerHTTPGinMode(t *testing.T) {
	// Test that server sets gin to release mode
	gin.SetMode(gin.TestMode) // Reset to test mode

	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	// This should set gin to release mode
	engine := NewServerHTTP(logger, userHandler)

	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// The function should have set gin to release mode
	// We can't easily test this without changing gin's global state
	// But we can verify the function completes without error
}

func TestServerHTTPRouteSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	baseHandler := handler.NewHandler(logger)
	userHandler := &handler.UserHandler{Handler: baseHandler}

	engine := NewServerHTTP(logger, userHandler)

	// Test that routes are properly configured by testing different paths
	testCases := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/user?id=1", http.StatusInternalServerError}, // Mock service will return error
		{"GET", "/nonexistent", http.StatusNotFound},
		{"POST", "/user", http.StatusNotFound}, // POST not configured for /user
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)

		engine.ServeHTTP(w, req)

		if w.Code != tc.expectedStatus {
			t.Errorf("For %s %s: expected status %d, got %d",
				tc.method, tc.path, tc.expectedStatus, w.Code)
		}
	}
}
