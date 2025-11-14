package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type TeamService struct {
	teamRepository interfaces.TeamRepository
}

func NewTeamService(repository interfaces.TeamRepository) TeamService {
	return TeamService{teamRepository: repository}
}

func (s *TeamService) Create(ctx context.Context, team_name string, members []entities.TeamMember) (*entities.Team, error) {
	team := entities.NewTeam(team_name, members)
	s.teamRepository.Create(ctx, team)
	return team, nil
}

func (s *TeamService) GetByID(ctx context.Context, team_name string, members []entities.TeamMember) (*entities.Team, error) {
	team := entities.NewTeam(team_name, members)
	s.teamRepository.Create(ctx, team)
	return team, nil
}
