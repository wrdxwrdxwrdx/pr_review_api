package interfaces

import (
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type UserRepository interface {
	repository.Repository[entities.User]

	SetIsActive(userId string, isActive bool) (*entities.User, error)
	GetUserTeam(userId string) (*[]string, error)
	GetById(userId string) (*entities.User, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
	// UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
}
