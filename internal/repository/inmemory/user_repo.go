package inmemory

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type UserRepository struct{}

func NewUserRepository() interfaces.UserRepository {
	return &UserRepository{}
}

func (r UserRepository) Create(ctx context.Context, entity *entities.User) error {
	return nil
}
