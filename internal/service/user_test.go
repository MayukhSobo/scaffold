package service

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/MayukhSobo/scaffold/internal/repository"
	"github.com/MayukhSobo/scaffold/pkg/log"
	_ "github.com/go-sql-driver/mysql"
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
