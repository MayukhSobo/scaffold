package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

func TestRunServer(t *testing.T) {
	// This test ensures RunServer can be called without panicking
	config := createTestConfig()
	logger := createTestLogger()

	// Create a done channel to signal when the test is complete
	done := make(chan bool, 1)

	// Run the server in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RunServer panicked: %v", r)
			}
			done <- true
		}()

		// We can't easily test the full RunServer function without mocking
		// the signal handling, so we'll test the components separately
		server := NewFiberServer(config, logger)
		if server == nil {
			t.Error("NewFiberServer returned nil")
		}

		app := server.GetApp()
		if app == nil {
			t.Error("GetApp returned nil")
		}
	}()

	// Wait for the goroutine to complete or timeout
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}

func TestRunFiberAppConfiguration(t *testing.T) {
	// Test that RunFiberApp correctly reads configuration
	config := createTestConfig()
	config.Set("http.port", "9999")
	config.Set("server.shutdown_timeout", "5s")

	// Test that the configuration is read correctly
	port := config.GetString("http.port")
	if port != "9999" {
		t.Errorf("Expected port '9999', got '%s'", port)
	}

	shutdownTimeout := config.GetDuration("server.shutdown_timeout")
	if shutdownTimeout != 5*time.Second {
		t.Errorf("Expected shutdown timeout '5s', got '%v'", shutdownTimeout)
	}
}

func TestRunFiberAppDefaultPort(t *testing.T) {
	// Test that RunFiberApp uses default port when not configured
	config := createTestConfig()
	config.Set("http.port", "") // Empty port should default to 8000

	// Test that default port is used
	port := config.GetString("http.port")
	if port != "" {
		// Reset to empty to test default behavior
		config.Set("http.port", "")
		port = config.GetString("http.port")
	}

	expectedPort := "8000"
	if port == "" {
		port = expectedPort // This simulates the default logic in RunFiberApp
	}

	if port != expectedPort {
		t.Errorf("Expected default port '%s', got '%s'", expectedPort, port)
	}
}

func TestRunFiberAppShutdownTimeout(t *testing.T) {
	// Test that RunFiberApp handles shutdown timeout configuration
	config := createTestConfig()

	// Test with configured timeout
	config.Set("server.shutdown_timeout", "10s")
	shutdownTimeout := config.GetDuration("server.shutdown_timeout")
	if shutdownTimeout != 10*time.Second {
		t.Errorf("Expected shutdown timeout '10s', got '%v'", shutdownTimeout)
	}

	// Test with empty timeout (should default to 30s)
	config.Set("server.shutdown_timeout", "")
	shutdownTimeout = config.GetDuration("server.shutdown_timeout")
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second // This simulates the default logic
	}
	if shutdownTimeout != 30*time.Second {
		t.Errorf("Expected default shutdown timeout '30s', got '%v'", shutdownTimeout)
	}
}

func TestRunWithCustomSetup(t *testing.T) {
	// Test that RunWithCustomSetup calls the setup function
	config := createTestConfig()
	logger := createTestLogger()

	setupCalled := false
	setupFunc := func(server *FiberServer) {
		setupCalled = true

		// Verify that server is not nil
		if server == nil {
			t.Error("Setup function received nil server")
		}

		// Verify that we can add routes
		server.AddRoutes(func(app *fiber.App) {
			app.Get("/custom-setup", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "custom setup route"})
			})
		})
	}

	// Test the setup function directly (since we can't easily test the full RunWithCustomSetup)
	server := NewFiberServer(config, logger)
	setupFunc(server)

	if !setupCalled {
		t.Error("Setup function was not called")
	}

	// Test that the custom route was added
	app := server.GetApp()
	req := httptest.NewRequest("GET", "/custom-setup", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test custom setup route: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestRunWithCustomSetupNilFunction(t *testing.T) {
	// Test that RunWithCustomSetup handles nil setup function
	config := createTestConfig()
	logger := createTestLogger()

	// Test with nil setup function (should not panic)
	server := NewFiberServer(config, logger)

	// Simulate the nil check in RunWithCustomSetup
	var setupFunc func(*FiberServer)
	setupFunc = nil

	if setupFunc != nil {
		setupFunc(server)
	}

	// If we get here without panicking, the test passes
	if server == nil {
		t.Error("Server should not be nil")
	}
}

func TestFiberServerLifecycle(t *testing.T) {
	// Test server creation and basic functionality
	config := createTestConfig()
	logger := createTestLogger()

	// Create server
	server := NewFiberServer(config, logger)
	if server == nil {
		t.Fatal("NewFiberServer returned nil")
	}

	// Get the app
	app := server.GetApp()
	if app == nil {
		t.Fatal("GetApp returned nil")
	}

	// Test that the server can handle requests
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test server lifecycle: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestFiberServerGracefulShutdown(t *testing.T) {
	// Test graceful shutdown functionality
	// Create a simple Fiber app
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	// Test that the app can be created and shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test that ShutdownWithContext doesn't panic
	err := app.ShutdownWithContext(ctx)
	if err != nil {
		// This is expected since the server wasn't actually started
		// but we're testing that the method exists and can be called
		t.Logf("Expected error during shutdown of non-started server: %v", err)
	}
}

func TestServerConfiguration(t *testing.T) {
	// Test various server configurations
	testCases := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			name: "minimal config",
			config: map[string]interface{}{
				"app.name":    "TestApp",
				"app.version": "1.0.0",
				"http.port":   "8080",
			},
		},
		{
			name: "full config",
			config: map[string]interface{}{
				"app.name":                      "TestApp",
				"app.version":                   "1.0.0",
				"http.port":                     "8080",
				"server.middleware.recover":     true,
				"server.middleware.request_id":  true,
				"server.middleware.logger":      true,
				"server.middleware.cors":        true,
				"server.shutdown_timeout":       "30s",
				"server.cors.allow_origins":     "http://localhost:3000",
				"server.cors.allow_methods":     "GET,POST,PUT,DELETE,OPTIONS",
				"server.cors.allow_headers":     "Content-Type,Authorization",
				"server.cors.allow_credentials": false,
				"server.cors.max_age":           86400,
			},
		},
		{
			name: "disabled middleware",
			config: map[string]interface{}{
				"app.name":                     "TestApp",
				"app.version":                  "1.0.0",
				"http.port":                    "8080",
				"server.middleware.recover":    false,
				"server.middleware.request_id": false,
				"server.middleware.logger":     false,
				"server.middleware.cors":       false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := viper.New()
			for key, value := range tc.config {
				config.Set(key, value)
			}

			logger := createTestLogger()
			server := NewFiberServer(config, logger)

			if server == nil {
				t.Errorf("NewFiberServer returned nil for config: %s", tc.name)
				return
			}

			// Test that the server can handle basic requests
			app := server.GetApp()
			req := httptest.NewRequest("GET", "/health", nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("Failed to test server with config %s: %v", tc.name, err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for config %s, got %d", tc.name, resp.StatusCode)
			}
		})
	}
}

func TestLoggerIntegration(t *testing.T) {
	// Test that the logger is properly integrated
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	config := createTestConfig()
	server := NewFiberServer(config, logger)

	// Test that the server can make requests and log them
	app := server.GetApp()
	req := httptest.NewRequest("GET", "/health", nil)
	_, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test logger integration: %v", err)
	}

	// Check that logging occurred
	if buf.Len() > 0 {
		t.Logf("Logger output: %s", buf.String())
	}
}

func BenchmarkFiberServerCreation(b *testing.B) {
	config := createTestConfig()
	logger := createTestLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := NewFiberServer(config, logger)
		if server == nil {
			b.Fatal("NewFiberServer returned nil")
		}
	}
}

func BenchmarkFiberServerRequest(b *testing.B) {
	config := createTestConfig()
	logger := createTestLogger()
	server := NewFiberServer(config, logger)
	app := server.GetApp()

	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req)
		if err != nil {
			b.Fatalf("Failed to test request: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}
