package service

import (
	"bytes"
	"github.com/MayukhSobo/scaffold/internal/repository"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"testing"

	"gorm.io/gorm"
)

func TestNewUserService(t *testing.T) {
	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	db := &gorm.DB{}
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	userService := NewUserService(service, userRepo)

	if userService == nil {
		t.Error("NewUserService() returned nil")
	}

	// Test that it implements UserService interface
	var _ UserService = userService
}

func TestUserServiceGetUserById(t *testing.T) {
	// Create test dependencies
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)

	service := NewService(logger)

	db := &gorm.DB{}
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	userService := NewUserService(service, userRepo)

	// Test GetUserById method
	user, err := userService.GetUserById(123)

	if err != nil {
		t.Errorf("GetUserById() returned error: %v", err)
	}

	if user == nil {
		t.Error("GetUserById() returned nil user")
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
	var serviceInterface UserService = userService
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
