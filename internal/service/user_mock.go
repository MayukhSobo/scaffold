package service

import (
	"github.com/MayukhSobo/scaffold/internal/model"
	"gorm.io/gorm"
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
	return &model.User{
		Model:    gorm.Model{ID: uint(id)},
		Username: "mock.user",
	}, nil
}
