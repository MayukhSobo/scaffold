package service

import (
	"gorm.io/gorm"

	"github.com/MayukhSobo/scaffold/internal/model"
)

// MockUserService is a mock implementation of UserService for testing.
type MockUserService struct {
	GetUserByIdFunc func(id int64) (*model.User, error)
}

// GetUserById implements the UserService interface for the mock.
func (m *MockUserService) GetUserById(id int64) (*model.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(id)
	}
	// Default mock behavior
	// Safely convert int64 to uint to prevent overflow
	var userID uint
	if id < 0 || id > int64(^uint(0)>>1) {
		userID = 1 // Default safe ID for mock
	} else {
		userID = uint(id)
	}

	return &model.User{
		Model:    gorm.Model{ID: userID},
		Username: "mock.user",
	}, nil
}
