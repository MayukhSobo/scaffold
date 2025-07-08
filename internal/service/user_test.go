package service

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/MayukhSobo/scaffold/internal/repository"
	"github.com/MayukhSobo/scaffold/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// setupTests initializes dependencies for testing.
func setupTests(t *testing.T) (UserService, *sql.DB) {
	// Use a in-memory SQLite database for testing
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/scaffold?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	queries := repository.New(db)
	userRepo := repository.NewUserRepository(queries)
	baseService := NewService(logger)
	userService := NewUserService(baseService, userRepo)

	return userService, db
}

func TestNewUserService(t *testing.T) {
	userService, _ := setupTests(t)
	if userService == nil {
		t.Error("NewUserService() returned nil")
	}
}

func TestUserServiceGetUserById(t *testing.T) {
	userService, db := setupTests(t)

	// Insert a test user
	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash, status, role) VALUES (?, ?, ?, ?, ?, ?)", 1, "testuser", "test@example.com", "hash", "active", "user")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

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

func TestUserServiceInterface(t *testing.T) {
	// Test that userService implements UserService interface
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	db := &gorm.DB{}
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	userService := NewUserService(service, userRepo)

	// Test interface compliance
	var serviceInterface = userService
	if serviceInterface == nil {
		t.Error("userService does not implement UserService interface")
	}

	// Test that methods exist
	_, err := serviceInterface.GetUserById(1)
	if err != nil {
		t.Errorf("Interface method GetUserById failed: %v", err)
	}
}

func TestUserServiceEmbedding(t *testing.T) {
	// Test that userService can be created successfully
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	db := &gorm.DB{}
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	userService := NewUserService(service, userRepo)

	// Test that it was created successfully
	if userService == nil {
		t.Error("userService should not be nil")
	}
}

func TestUserServiceDependencyInjection(t *testing.T) {
	// Test that dependency injection works correctly
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	db := &gorm.DB{}
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	userService := NewUserService(service, userRepo)

	// Test that the service can execute methods
	_, err := userService.GetUserById(1)
	if err != nil {
		t.Errorf("Dependency injection failed: %v", err)
	}
}
