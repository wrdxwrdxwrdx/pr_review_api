package interfaces

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type UserRepository interface {
	repository.Repository[entities.User]

	SetIsActive(ctx context.Context, userId string, isActive bool) (*entities.User, error)
	GetUserTeam(ctx context.Context, teamName string) ([]string, error)
	GetById(ctx context.Context, userId string) (*entities.User, error)
}
