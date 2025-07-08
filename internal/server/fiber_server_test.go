package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

func createTestConfig() *viper.Viper {
	config := viper.New()
	config.SetDefault("app.name", "TestApp")
	config.SetDefault("app.version", "1.0.0")
	config.SetDefault("env", "test")
	config.SetDefault("http.port", "8080")

	// Server middleware settings
	config.SetDefault("server.middleware.recover", true)
	config.SetDefault("server.middleware.request_id", true)
	config.SetDefault("server.middleware.logger", true)
	config.SetDefault("server.middleware.logger_format", "[${time}] ${status} - ${method} ${path}\n")
	config.SetDefault("server.middleware.cors", true)

	// CORS settings
	config.SetDefault("server.cors.allow_origins", "http://localhost:3000")
	config.SetDefault("server.cors.allow_methods", "GET,POST,PUT,DELETE,OPTIONS")
	config.SetDefault("server.cors.allow_headers", "Content-Type,Authorization")
	config.SetDefault("server.cors.allow_credentials", false)
	config.SetDefault("server.cors.max_age", 86400)

	return config
}

func createTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}

func TestNewFiberServer(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	if server == nil {
		t.Fatal("NewFiberServer should not return nil")
	}

	if server.app == nil {
		t.Error("FiberServer app should not be nil")
	}

	if server.config == nil {
		t.Error("FiberServer config should not be nil")
	}

	if server.logger == nil {
		t.Error("FiberServer logger should not be nil")
	}
}

func TestFiberServerGetApp(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	if app == nil {
		t.Error("GetApp should not return nil")
	}

	// Test that the app is a proper Fiber app
	if app != server.app {
		t.Error("GetApp should return the same app instance")
	}
}

func TestFiberServerHealthEndpoint(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test health endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if response["env"] != "test" {
		t.Errorf("Expected env 'test', got %v", response["env"])
	}
}

func TestFiberServerPingEndpoint(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	// Create a test request
	req := httptest.NewRequest("GET", "/ping", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test ping endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["message"] != "pong" {
		t.Errorf("Expected message 'pong', got %v", response["message"])
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestFiberServerRootEndpoint(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test root endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	expectedMessage := "Welcome to TestApp"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got %v", expectedMessage, response["message"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %v", response["version"])
	}

	if response["status"] != "running" {
		t.Errorf("Expected status 'running', got %v", response["status"])
	}
}

func TestFiberServerAddRoutes(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	// Add custom routes
	server.AddRoutes(func(app *fiber.App) {
		app.Get("/custom", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "custom route"})
		})
	})

	app := server.GetApp()

	// Test the custom route
	req := httptest.NewRequest("GET", "/custom", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test custom route: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["message"] != "custom route" {
		t.Errorf("Expected message 'custom route', got %v", response["message"])
	}
}

func TestFiberServerAddGroup(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	// Add a route group
	server.AddGroup("/api/v1", func(router fiber.Router) {
		router.Get("/users", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"users": []string{"user1", "user2"}})
		})
	})

	app := server.GetApp()

	// Test the grouped route
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test grouped route: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	users, ok := response["users"].([]interface{})
	if !ok {
		t.Error("Expected users to be an array")
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestFiberServerAddMiddleware(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	// Add custom middleware
	middlewareCalled := false
	server.AddMiddleware(func(c *fiber.Ctx) error {
		middlewareCalled = true
		return c.Next()
	})

	// Add a custom route after middleware to test it
	server.AddRoutes(func(app *fiber.App) {
		app.Get("/test-middleware", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "middleware test"})
		})
	})

	app := server.GetApp()

	// Test that middleware is called on the new route
	req := httptest.NewRequest("GET", "/test-middleware", nil)
	_, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test middleware: %v", err)
	}

	if !middlewareCalled {
		t.Error("Custom middleware was not called")
	}
}

func TestFiberServerErrorHandler(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	// Add a route that returns an error
	server.AddRoutes(func(app *fiber.App) {
		app.Get("/error", func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusBadRequest, "Test error")
		})
	})

	app := server.GetApp()

	// Test error handling
	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test error handler: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["error"] != true {
		t.Errorf("Expected error to be true, got %v", response["error"])
	}

	if response["message"] != "Test error" {
		t.Errorf("Expected message 'Test error', got %v", response["message"])
	}

	if response["code"] != float64(400) {
		t.Errorf("Expected code 400, got %v", response["code"])
	}
}

func TestFiberServerWithDisabledMiddleware(t *testing.T) {
	config := createTestConfig()
	// Disable all middleware
	config.Set("server.middleware.recover", false)
	config.Set("server.middleware.request_id", false)
	config.Set("server.middleware.logger", false)
	config.Set("server.middleware.cors", false)

	logger := createTestLogger()

	server := NewFiberServer(config, logger)

	if server == nil {
		t.Fatal("NewFiberServer should not return nil even with disabled middleware")
	}

	// Test that basic endpoints still work
	app := server.GetApp()
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test with disabled middleware: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestFiberServerCORSConfiguration(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	// Test CORS preflight request
	req := httptest.NewRequest("OPTIONS", "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test CORS: %v", err)
	}

	// Check CORS headers
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be 'http://localhost:3000', got '%s'", allowOrigin)
	}
}

func TestFiberServerConfiguration(t *testing.T) {
	config := createTestConfig()
	logger := createTestLogger()

	server := NewFiberServer(config, logger)
	app := server.GetApp()

	// Test that server configuration is applied
	if app == nil {
		t.Fatal("Fiber app should not be nil")
	}

	// Test that we can make requests (indicating the server is configured)
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test server configuration: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
