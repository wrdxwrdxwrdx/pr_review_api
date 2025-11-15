package interfaces

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type PrRepository interface {
	repository.Repository[entities.PullRequest]

	GetByID(ctx context.Context, prId string) (*entities.PullRequest, error)
	Merge(ctx context.Context, prId string) (*entities.PullRequest, error)
	// GetByEmail(ctx context.Context, email string) (*entities.User, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
	// UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
}
