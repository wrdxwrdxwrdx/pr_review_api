package inmemory

import (
	"context"
	"fmt"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
	"time"
)

type PrRepository struct {
	prs map[string]*entities.PullRequest
}

func NewPrRepository() interfaces.PrRepository {
	return &PrRepository{
		prs: make(map[string]*entities.PullRequest),
	}
}

func (r *PrRepository) Create(ctx context.Context, entity *entities.PullRequest) error {
	fmt.Println("Create Pr inmemory")
	r.prs[entity.PullRequestId] = entity
	return nil
}

func (r *PrRepository) GetByID(ctx context.Context, prId string) (*entities.PullRequest, error) {
	fmt.Println("Get Pr inmemory")
	return r.prs[prId], nil
}

func (r *PrRepository) Merge(ctx context.Context, prId string) (*entities.PullRequest, error) {
	pr, err := r.GetByID(ctx, prId)

	if pr.Status == entities.StatusMerged {
		return pr, nil
	}

	if err != nil {
		return nil, err
	}

	mergedAt := time.Now()
	pr.MergedAt = &mergedAt
	return pr, nil
}

func (r *PrRepository) Reassign(ctx context.Context, pullRequestId string, newAssignedReviewers []string) error {
	pr, err := r.GetByID(ctx, pullRequestId)

	if err != nil {
		return err
	}

	pr.AssignedReviewers = newAssignedReviewers
	return nil
}

func (r *PrRepository) GetReview(ctx context.Context, authorId string) ([]*entities.PullRequestShort, error) {
	var reviewPrs []*entities.PullRequestShort

	for _, pr := range r.prs {
		if pr.AuthorId == authorId && pr.Status == entities.StatusOpen {
			pullRequestShort := entities.NewPullRequestShortFromPr(*pr)
			reviewPrs = append(reviewPrs, pullRequestShort)
		}
	}

	return reviewPrs, nil
}
