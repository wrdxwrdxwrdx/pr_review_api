package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
	"time"
)

type PrService struct {
	PrRepository   interfaces.PrRepository
	UserRepository interfaces.UserRepository
}

func NewPrService(prRepository interfaces.PrRepository, userRepository interfaces.UserRepository) PrService {
	return PrService{PrRepository: prRepository, UserRepository: userRepository}
}

func (s *PrService) Create(ctx context.Context, pullRequestId string, pullRequestName string, authorId string) (*entities.PullRequest, error) {
	status := entities.StatusOpen
	teamMembers, err := s.UserRepository.GetUserTeam(authorId)

	if err != nil {
		return nil, err
	}

	var assignedReviewers []string
	for i := 0; i < 2; i++ {
		if i < len(*teamMembers) {
			assignedReviewers = append(assignedReviewers, (*teamMembers)[i])
		}
	}

	createdAt := time.Now()

	pr := entities.NewPullRequest(pullRequestId, pullRequestName, authorId, status, assignedReviewers, &createdAt)
	s.PrRepository.Create(ctx, pr)
	return pr, nil
}

func (s *PrService) Merge(ctx context.Context, pullRequestId string) (*entities.PullRequest, error) {
	return s.PrRepository.Merge(ctx, pullRequestId)
}
