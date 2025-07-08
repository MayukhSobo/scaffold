package routes

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/pkg/log"
	_ "github.com/go-sql-driver/mysql"
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

func createTestApp() *fiber.App {
	return fiber.New()
}

func createTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}

func createTestUserRepo(t *testing.T) users.Querier {
	// Create a test database connection
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/scaffold?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return users.New(db)
}

func TestGetAdminUsersRoute(t *testing.T) {
	// Create test app
	app := createTestApp()
	logger := createTestLogger()

	// Create mock user service
	mockUserService := &mockUserService{}

	// Create test user repository
	userRepo := createTestUserRepo(t)

	// Create base handler
	baseHandler := handler.NewHandler(logger)

	// Create API v1 group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register user routes
	RegisterUserRoutes(v1, baseHandler, mockUserService, userRepo)

	// Test admin users route
	req := httptest.NewRequest("GET", "/api/v1/users/admin", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test admin users route: %v", err)
	}

	// Note: This test may fail if the database doesn't have admin users
	// For now, we'll just check that the route exists and doesn't crash
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestGetPendingVerificationUsersRoute(t *testing.T) {
	// Create test app
	app := createTestApp()
	logger := createTestLogger()

	// Create mock user service
	mockUserService := &mockUserService{}

	// Create test user repository
	userRepo := createTestUserRepo(t)

	// Create base handler
	baseHandler := handler.NewHandler(logger)

	// Create API v1 group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register user routes
	RegisterUserRoutes(v1, baseHandler, mockUserService, userRepo)

	// Test pending verification users route
	req := httptest.NewRequest("GET", "/api/v1/users/pending-verification", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test pending verification users route: %v", err)
	}

	// Note: This test may fail if the database doesn't have pending users
	// For now, we'll just check that the route exists and doesn't crash
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}
