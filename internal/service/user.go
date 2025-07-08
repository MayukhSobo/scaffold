package service

import (
	"context"

	"github.com/MayukhSobo/scaffold/internal/repository/users"
)

type UserService interface {
	GetUserById(ctx context.Context, id int64) (users.User, error)
}

type userService struct {
	*Service
	userRepository users.Querier
}

func NewUserService(service *Service, userRepository users.Querier) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
	}
}

func (s *userService) GetUserById(ctx context.Context, id int64) (users.User, error) {
	return s.userRepository.GetUser(ctx, uint64(id))
}
