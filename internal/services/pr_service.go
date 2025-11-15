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

func (s *PrService) Reassign(ctx context.Context, pullRequestId string, oldReviewerId string) (*entities.PullRequest, error) {
	pr, err := s.PrRepository.GetByID(ctx, pullRequestId)
	if err != nil {
		return nil, err
	}

	newAssignedReviewers := pr.AssignedReviewers
	team, err := s.UserRepository.GetUserTeam(pr.AuthorId)

	if err != nil {
		return nil, err
	}

	var newReviewer string
	var reviewerOne, reviewerTwo string
	if len(newAssignedReviewers) > 0 {
		reviewerOne = newAssignedReviewers[0]
	}
	if len(newAssignedReviewers) > 1 {
		reviewerTwo = newAssignedReviewers[1]
	}

	for _, user := range *team {
		if user != oldReviewerId && user != reviewerOne && user != reviewerTwo {
			newReviewer = user
			break
		}
	}

	for i, user := range newAssignedReviewers {
		if user == oldReviewerId {
			newAssignedReviewers[i] = newReviewer
			if newReviewer == "" {
				newAssignedReviewers = append(newAssignedReviewers[:i], newAssignedReviewers[i+1:]...)
			}

		}
	}

	return s.PrRepository.Reassign(ctx, pullRequestId, newAssignedReviewers)
}
