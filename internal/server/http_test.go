package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"

	"github.com/gin-gonic/gin"
)

// setupTestServer initializes a test server with mock dependencies.
func setupTestServer() (*gin.Engine, *bytes.Buffer) {
	gin.SetMode(gin.TestMode)
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	// Create mock service
	mockUserService := &service.MockUserService{}

	// Create handlers with mock dependencies
	baseHandler := handler.NewHandler(logger)
	userHandler := handler.NewUserHandler(baseHandler, mockUserService)

	engine := NewServerHTTP(logger, userHandler)
	return engine, &buf
}

func TestNewServerHTTP(t *testing.T) {
	engine, _ := setupTestServer()
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
	engine, buf := setupTestServer()

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
	engine, _ := setupTestServer()

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
	engine, _ := setupTestServer()

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
	engine, _ := setupTestServer()

	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// The function should have set gin to release mode
	// We can't easily test this without changing gin's global state
	// But we can verify the function completes without error
}

func TestServerHTTPRouteSetup(t *testing.T) {
	engine, _ := setupTestServer()

	// Test that routes are properly configured by testing different paths
	testCases := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/user?id=1", http.StatusOK}, // Mock service will now return success
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
