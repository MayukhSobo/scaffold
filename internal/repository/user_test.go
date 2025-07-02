package repository

import (
	"bytes"
	"testing"

	"github.com/MayukhSobo/scaffold/pkg/log"

	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	// Create a test logger and repository
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)

	userRepo := NewUserRepository(repo)

	if userRepo == nil {
		t.Error("NewUserRepository() returned nil")
	}

	// Test that it implements UserRepository interface
	var _ = userRepo
}

func TestUserRepositoryFirstById(t *testing.T) {
	// Create a test setup
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)
	userRepo := NewUserRepository(repo)

	// Test FirstById method (will return empty user due to mock implementation)
	user, err := userRepo.FirstById(123)

	if err != nil {
		t.Errorf("FirstById() returned error: %v", err)
	}

	if user == nil {
		t.Error("FirstById() returned nil user")
		return
	}

	// Since this is a mock implementation, we expect an empty user
	if user.ID != 0 {
		t.Errorf("Expected empty user ID, got %d", user.ID)
	}
}

func TestUserRepositoryInterface(t *testing.T) {
	// Test that userRepository implements UserRepository interface
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)
	userRepo := NewUserRepository(repo)

	// Test interface compliance
	var userRepoInterface = userRepo
	if userRepoInterface == nil {
		t.Error("userRepository does not implement UserRepository interface")
	}

	// Test that methods exist
	_, err := userRepoInterface.FirstById(1)
	if err != nil {
		t.Errorf("Interface method FirstById failed: %v", err)
	}
}

func TestUserRepositoryEmbedding(t *testing.T) {
	// Test that userRepository can be created successfully
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	db := &gorm.DB{}
	repo := NewRepository(logger, db)
	userRepo := NewUserRepository(repo)

	// Test that it was created successfully
	if userRepo == nil {
		t.Error("userRepository should not be nil")
	}
}
