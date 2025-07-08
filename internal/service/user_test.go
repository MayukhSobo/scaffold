package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

// mockUserRepository implements users.Querier for testing
type mockUserRepository struct {
	users []users.User
}

func (m *mockUserRepository) GetUser(ctx context.Context, id uint64) (users.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return users.User{}, nil // Return empty user if not found
}

func (m *mockUserRepository) GetUsers(ctx context.Context) ([]users.User, error) {
	return m.users, nil
}

func (m *mockUserRepository) GetAdminUsers(ctx context.Context) ([]users.User, error) {
	var adminUsers []users.User
	for _, user := range m.users {
		if user.Role == "admin" {
			adminUsers = append(adminUsers, user)
		}
	}
	return adminUsers, nil
}

func (m *mockUserRepository) GetPendingVerificationUsers(ctx context.Context) ([]users.User, error) {
	var pendingUsers []users.User
	for _, user := range m.users {
		if user.Status == "pending_verification" {
			pendingUsers = append(pendingUsers, user)
		}
	}
	return pendingUsers, nil
}

// setupTestsWithMock initializes dependencies for testing using mocks
func setupTestsWithMock(t *testing.T) (UserService, *mockUserRepository) {
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	// Create mock repository with test data
	mockRepo := &mockUserRepository{
		users: []users.User{
			{
				ID:           1,
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hash",
				Status:       "active",
				Role:         "user",
			},
			{
				ID:           2,
				Username:     "admin",
				Email:        "admin@example.com",
				PasswordHash: "hash",
				Status:       "active",
				Role:         "admin",
			},
			{
				ID:           3,
				Username:     "pending",
				Email:        "pending@example.com",
				PasswordHash: "hash",
				Status:       "pending_verification",
				Role:         "user",
			},
		},
	}

	baseService := NewService(logger)
	userService := NewUserService(baseService, mockRepo)

	return userService, mockRepo
}

func TestNewUserService(t *testing.T) {
	userService, _ := setupTestsWithMock(t)
	if userService == nil {
		t.Error("NewUserService() returned nil")
	}
}

func TestUserServiceGetUserById(t *testing.T) {
	userService, _ := setupTestsWithMock(t)

	user, err := userService.GetUserById(context.Background(), 1)
	if err != nil {
		t.Errorf("GetUserById() returned error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
}

func TestUserServiceGetAdminUsers(t *testing.T) {
	userService, _ := setupTestsWithMock(t)

	adminUsers, err := userService.GetAdminUsers(context.Background())
	if err != nil {
		t.Errorf("GetAdminUsers() returned error: %v", err)
	}

	if len(adminUsers) != 1 {
		t.Errorf("Expected 1 admin user, got %d", len(adminUsers))
	}

	if len(adminUsers) > 0 && adminUsers[0].Username != "admin" {
		t.Errorf("Expected admin username 'admin', got %s", adminUsers[0].Username)
	}
}

func TestUserServiceGetPendingVerificationUsers(t *testing.T) {
	userService, _ := setupTestsWithMock(t)

	pendingUsers, err := userService.GetPendingVerificationUsers(context.Background())
	if err != nil {
		t.Errorf("GetPendingVerificationUsers() returned error: %v", err)
	}

	if len(pendingUsers) != 1 {
		t.Errorf("Expected 1 pending user, got %d", len(pendingUsers))
	}

	if len(pendingUsers) > 0 && pendingUsers[0].Username != "pending" {
		t.Errorf("Expected pending username 'pending', got %s", pendingUsers[0].Username)
	}
}

func TestUserServiceGetUserByIdNotFound(t *testing.T) {
	userService, _ := setupTestsWithMock(t)

	user, err := userService.GetUserById(context.Background(), 999)
	if err != nil {
		t.Errorf("GetUserById() returned error: %v", err)
	}

	if user.ID != 0 {
		t.Errorf("Expected empty user (ID 0) for non-existent user, got ID %d", user.ID)
	}
}
