package container

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

func createTestConfig() *viper.Viper {
	conf := viper.New()
	conf.Set("db.mysql.host", "localhost")
	conf.Set("db.mysql.port", "3306")
	conf.Set("db.mysql.user", "test")
	conf.Set("db.mysql.password", "test")
	conf.Set("db.mysql.database", "test")
	return conf
}

func createTestLogger() log.Logger {
	var buf bytes.Buffer
	return log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
}

func TestNewTypedContainer(t *testing.T) {
	conf := createTestConfig()
	logger := createTestLogger()

	// We can't create a real database connection in tests,
	// so we'll test the container structure without database
	container := &TypedContainer{
		config: conf,
		logger: logger,
	}

	if container.GetConfig() != conf {
		t.Error("Container should return the correct config")
	}

	if container.GetLogger() != logger {
		t.Error("Container should return the correct logger")
	}
}

func TestTypedContainerGetters(t *testing.T) {
	conf := createTestConfig()
	logger := createTestLogger()

	container := &TypedContainer{
		config: conf,
		logger: logger,
	}

	// Test infrastructure getters
	if container.GetConfig() == nil {
		t.Error("GetConfig() should not return nil")
	}

	if container.GetLogger() == nil {
		t.Error("GetLogger() should not return nil")
	}
}

func TestAllServicesStruct(t *testing.T) {
	// Test that AllServices struct can hold service interfaces
	var buf bytes.Buffer
	logger := log.NewConsoleLoggerWithWriter(log.InfoLevel, &buf, false)
	baseService := service.NewService(logger)

	// Create a mock user service for testing
	// In real usage, this would come from the container
	services := &AllServices{
		User: nil, // Would be populated by container
	}

	// Test struct exists and can be used
	if services == nil {
		t.Error("AllServices struct should be creatable")
	}

	// Test with actual service
	// This would normally be done by the container
	_ = baseService // Use the baseService to avoid unused variable
}

func TestAllRepositoriesStruct(t *testing.T) {
	// Test that AllRepositories struct can hold repository interfaces
	repositories := &AllRepositories{
		User: nil, // Would be populated by container
	}

	// Test struct exists and can be used
	if repositories == nil {
		t.Error("AllRepositories struct should be creatable")
	}
}

func TestContainerInterface(t *testing.T) {
	// Test that the container provides the expected interface
	conf := createTestConfig()
	logger := createTestLogger()

	container := &TypedContainer{
		config: conf,
		logger: logger,
	}

	// Test that all getter methods exist and return proper types
	config := container.GetConfig()
	if config == nil {
		t.Error("GetConfig should return a config")
	}

	containerLogger := container.GetLogger()
	if containerLogger == nil {
		t.Error("GetLogger should return a logger")
	}

	// Test service getters (would return nil without proper initialization)
	userService := container.GetUserService()
	_ = userService // May be nil in test, that's ok

	// Test repository getters (would return nil without proper initialization)
	userRepo := container.GetUserRepository()
	_ = userRepo // May be nil in test, that's ok
}

func TestGetAllServices(t *testing.T) {
	container := &TypedContainer{}

	allServices := container.GetAllServices()
	if allServices == nil {
		t.Error("GetAllServices should return a non-nil AllServices struct")
	}

	// Test that the struct has the expected fields
	_ = allServices.User // Should exist
}

func TestGetAllRepositories(t *testing.T) {
	container := &TypedContainer{}

	allRepos := container.GetAllRepositories()
	if allRepos == nil {
		t.Error("GetAllRepositories should return a non-nil AllRepositories struct")
	}

	// Test that the struct has the expected fields
	_ = allRepos.User // Should exist
}

// Mock implementations for testing

type mockUserRepository struct{}

func (m *mockUserRepository) GetUser(ctx context.Context, id uint64) (users.User, error) {
	return users.User{ID: id, Username: "test"}, nil
}

func (m *mockUserRepository) GetUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
	}, nil
}

func (m *mockUserRepository) GetAdminUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{{ID: 1, Username: "admin"}}, nil
}

func (m *mockUserRepository) GetPendingVerificationUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{{ID: 2, Username: "pending"}}, nil
}

func TestContainerWithMockDependencies(t *testing.T) {
	// This demonstrates how the container can work with mock dependencies for testing
	conf := createTestConfig()
	logger := createTestLogger()

	container := &TypedContainer{
		config:         conf,
		logger:         logger,
		userRepository: &mockUserRepository{},
	}

	// Test that we can get the mock repository
	userRepo := container.GetUserRepository()
	if userRepo == nil {
		t.Error("Container should return the user repository")
	}

	// Test that the mock works
	ctx := context.Background()
	users, err := userRepo.GetAdminUsers(ctx)
	if err != nil {
		t.Errorf("Mock repository should not return error: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0].Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", users[0].Username)
	}
}

// Example test showing how container makes testing easier
func TestContainerDrivenHandler(t *testing.T) {
	// Setup container with mocks
	container := &TypedContainer{
		config:         createTestConfig(),
		logger:         createTestLogger(),
		userRepository: &mockUserRepository{},
	}

	// Create service with mocked dependencies
	baseService := service.NewService(container.GetLogger())
	container.userService = service.NewUserService(baseService, container.GetUserRepository())

	// Test that services work through container
	userService := container.GetUserService()
	if userService == nil {
		t.Error("Container should provide user service")
	}

	// This demonstrates how handlers would use the container in tests
	allServices := container.GetAllServices()
	if allServices.User == nil {
		t.Error("All services should include user service")
	}
}
