package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

func createTestApp() *fiber.App {
	return fiber.New()
}

func createTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}

func TestGetAdminUsersRoute(t *testing.T) {
	// Create test app
	app := createTestApp()
	logger := createTestLogger()

	// Create mock user service
	mockUserService := &service.MockUserService{}

	// Create base handler
	baseHandler := handler.NewHandler(logger)

	// Create API v1 group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register user routes
	RegisterUserRoutes(v1, baseHandler, mockUserService)

	// Test admin users route
	req := httptest.NewRequest("GET", "/api/v1/users/admin", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test admin users route: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check response structure
	if response["code"] != float64(0) {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got %v", response["message"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected data to be an object")
	}

	users, ok := data["users"].([]interface{})
	if !ok {
		t.Error("Expected users to be an array")
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 admin users, got %d", len(users))
	}

	// Check first user has admin role
	firstUser, ok := users[0].(map[string]interface{})
	if !ok {
		t.Error("Expected first user to be an object")
	}

	if firstUser["role"] != "admin" {
		t.Errorf("Expected first user to have admin role, got %v", firstUser["role"])
	}
}

func TestGetPendingVerificationUsersRoute(t *testing.T) {
	// Create test app
	app := createTestApp()
	logger := createTestLogger()

	// Create mock user service
	mockUserService := &service.MockUserService{}

	// Create base handler
	baseHandler := handler.NewHandler(logger)

	// Create API v1 group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register user routes
	RegisterUserRoutes(v1, baseHandler, mockUserService)

	// Test pending verification users route
	req := httptest.NewRequest("GET", "/api/v1/users/pending-verification", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test pending verification users route: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check response structure
	if response["code"] != float64(0) {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got %v", response["message"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected data to be an object")
	}

	users, ok := data["users"].([]interface{})
	if !ok {
		t.Error("Expected users to be an array")
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 pending verification users, got %d", len(users))
	}

	// Check first user structure
	firstUser, ok := users[0].(map[string]interface{})
	if !ok {
		t.Error("Expected first user to be an object")
	}

	if firstUser["status"] != "pending_verification" {
		t.Errorf("Expected status 'pending_verification', got %v", firstUser["status"])
	}

	if firstUser["email"] == nil {
		t.Error("Expected user to have email field")
	}

	if firstUser["verification_token"] == nil {
		t.Error("Expected user to have verification_token field")
	}
}
