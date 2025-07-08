package routes

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

// mockUserService implements service.UserService for testing
type mockUserService struct{}

func (m *mockUserService) GetUserById(ctx context.Context, id int64) (users.User, error) {
	return users.User{
		ID:       uint64(id),
		Username: "testuser",
		Email:    "test@example.com",
	}, nil
}

func (m *mockUserService) GetAdminUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{
		{
			ID:       1,
			Username: "admin",
			Email:    "admin@example.com",
			Role:     users.UsersRoleAdmin,
		},
		{
			ID:       2,
			Username: "superadmin",
			Email:    "superadmin@example.com",
			Role:     users.UsersRoleAdmin,
		},
	}, nil
}

func (m *mockUserService) GetPendingVerificationUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{
		{
			ID:       3,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   users.UsersStatusPendingVerification,
		},
		{
			ID:       4,
			Username: "user2",
			Email:    "user2@example.com",
			Status:   users.UsersStatusPendingVerification,
		},
	}, nil
}

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
	mockUserService := &mockUserService{}

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
}

func TestGetPendingVerificationUsersRoute(t *testing.T) {
	// Create test app
	app := createTestApp()
	logger := createTestLogger()

	// Create mock user service
	mockUserService := &mockUserService{}

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
}
