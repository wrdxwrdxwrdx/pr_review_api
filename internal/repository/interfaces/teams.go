package interfaces

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type TeamRepository interface {
	repository.Repository[entities.Team]

	GetByName(ctx context.Context, teamNameQuery string) (*entities.Team, error)
}
