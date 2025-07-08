package repository

import (
	"context"
)

type UserRepository interface {
	GetAdminUsers(ctx context.Context) ([]User, error)
	GetPendingVerificationUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uint64) (User, error)
}

type userRepository struct {
	*Queries
}

func NewUserRepository(q *Queries) UserRepository {
	return &userRepository{
		Queries: q,
	}
}

func (r *userRepository) GetAdminUsers(ctx context.Context) ([]User, error) {
	return r.Queries.GetAdminUsers(ctx)
}

func (r *userRepository) GetPendingVerificationUsers(ctx context.Context) ([]User, error) {
	return r.Queries.GetPendingVerificationUsers(ctx)
}

func (r *userRepository) GetUser(ctx context.Context, id uint64) (User, error) {
	return r.Queries.GetUser(ctx, id)
}
