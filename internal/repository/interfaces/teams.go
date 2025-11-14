package interfaces

import (
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type TeamRepository interface {
	repository.Repository[entities.Team]

	// GetByEmail(ctx context.Context, email string) (*entities.User, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
	// UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
}
