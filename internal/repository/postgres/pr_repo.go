package postgres

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/domain/errors"
	"pr_review_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PrRepository struct {
	db *pgxpool.Pool
}

func NewPRRepository(db *pgxpool.Pool) *PrRepository {
	return &PrRepository{db: db}
}

func (r *PrRepository) Create(ctx context.Context, entity *entities.PullRequest) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO pull_requests 
        (pull_request_id, pull_request_name, author_id, status, assigned_reviewers) 
        VALUES ($1, $2, $3, $4, $5)
    `, entity.PullRequestId, entity.PullRequestName, entity.AuthorId, entity.Status, entity.AssignedReviewers)

	if err != nil {
		if utils.IsUnique(err) {
			return errors.NewDomainError(errors.PRExists, "PR already exists")
		}
		return err
	}
	return nil
}

func (r *PrRepository) GetByID(ctx context.Context, prId string) (*entities.PullRequest, error) {
	var pr entities.PullRequest
	err := r.db.QueryRow(ctx, `
        SELECT pull_request_id, pull_request_name, author_id, status, 
               assigned_reviewers, created_at, merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1
    `, prId).Scan(
		&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status,
		&pr.AssignedReviewers, &pr.CreatedAt, &pr.MergedAt,
	)
	if err != nil {
		return nil, errors.NewDomainError(errors.NotFound, "PR not found")
	}
	return &pr, nil
}

func (r *PrRepository) Merge(ctx context.Context, prId string) (*entities.PullRequest, error) {
	pr, err := r.GetByID(ctx, prId)

	if err != nil {
		return nil, err
	}

	if pr.Status == entities.StatusMerged {
		return pr, err
	}

	var newPr entities.PullRequest

	err = r.db.QueryRow(ctx, `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = CURRENT_TIMESTAMP 
        WHERE pull_request_id = $1 
        RETURNING pull_request_id, pull_request_name, author_id, status, 
                  assigned_reviewers, created_at, merged_at
    `, prId).Scan(
		&newPr.PullRequestId, &newPr.PullRequestName, &newPr.AuthorId, &newPr.Status,
		&newPr.AssignedReviewers, &newPr.CreatedAt, &newPr.MergedAt,
	)
	if err != nil {
		return nil, errors.NewDomainError(errors.NotFound, "PR not found")
	}
	return &newPr, nil
}

func (r *PrRepository) Reassign(ctx context.Context, pullRequestId string, newAssignedReviewers []string) error {
	pr, error := r.GetByID(ctx, pullRequestId)

	if error != nil {
		return error
	}

	_, err := r.db.Exec(ctx, `
        UPDATE pull_requests 
        SET assigned_reviewers = $1, status = $2, merged_at = $3
        WHERE pull_request_id = $4
    `, newAssignedReviewers, pr.Status, pr.MergedAt, pr.PullRequestId)
	return err
}

func (r *PrRepository) GetReview(ctx context.Context, authorId string) ([]*entities.PullRequestShort, error) {
	rows, err := r.db.Query(ctx, `
        SELECT pull_request_id, pull_request_name, author_id, status
        FROM pull_requests 
        WHERE $1 = ANY(assigned_reviewers) AND status = 'OPEN'
    `, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*entities.PullRequestShort
	for rows.Next() {
		var pr entities.PullRequestShort
		err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}
