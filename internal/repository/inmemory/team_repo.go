package inmemory

import (
	"context"
	"fmt"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type TeamRepository struct {
	teams   []entities.Team
	members []entities.TeamMember
}

func NewTeamRepository() interfaces.TeamRepository {
	return &TeamRepository{}
}

func (r *TeamRepository) Create(ctx context.Context, entity *entities.Team) error {
	fmt.Println("Create team inmemory")
	r.teams = append(r.teams, *entity)
	for _, member := range entity.Members {
		fmt.Println("create member")
		r.members = append(r.members, member)
	}
	return nil
}
