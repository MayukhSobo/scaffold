package service

import (
	"context"

	"github.com/MayukhSobo/scaffold/internal/repository"
)

type UserService interface {
	GetUserById(ctx context.Context, id int64) (repository.User, error)
}

type userService struct {
	*Service
	userRepository repository.UserRepository
}

func NewUserService(service *Service, userRepository repository.UserRepository) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
	}
}

func (s *userService) GetUserById(ctx context.Context, id int64) (repository.User, error) {
	return s.userRepository.GetUser(ctx, uint64(id))
}
