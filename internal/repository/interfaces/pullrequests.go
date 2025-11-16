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
	Reassign(ctx context.Context, pullRequestId string, newAssignedReviewers []string) error
	GetReview(ctx context.Context, authorId string) ([]*entities.PullRequestShort, error)
}
