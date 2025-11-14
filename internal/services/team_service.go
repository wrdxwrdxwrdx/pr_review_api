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

func (s *TeamService) Create(ctx context.Context, teamName string, members []entities.TeamMember) (*entities.Team, error) {
	team := entities.NewTeam(teamName, members)
	s.teamRepository.Create(ctx, team)
	return team, nil
}

func (s *TeamService) GetByID(ctx context.Context, teamNameQuery string) (*entities.Team, error) {
	team, err := s.teamRepository.GetByID(ctx, teamNameQuery)
	return team, err
}
