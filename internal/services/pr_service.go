package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"
	"pr_review_api/internal/repository/interfaces"
	"slices"
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
	author, err := s.UserRepository.GetById(ctx, authorId)

	if err != nil {
		return nil, err
	}

	teamMembers, err := s.UserRepository.GetUserTeam(ctx, author.TeamName)

	if err != nil {
		return nil, err
	}

	var assignedReviewers []string
	i := 0
	for _, userId := range teamMembers {
		if i < 2 && userId != authorId {
			i += 1
			assignedReviewers = append(assignedReviewers, (teamMembers)[i])
		}
	}

	createdAt := time.Now()

	pr := entities.NewPullRequest(pullRequestId, pullRequestName, authorId, status, assignedReviewers, &createdAt)
	err = s.PrRepository.Create(ctx, pr)
	return pr, err
}

func (s *PrService) Merge(ctx context.Context, pullRequestId string) (*entities.PullRequest, error) {
	return s.PrRepository.Merge(ctx, pullRequestId)
}

func (s *PrService) Reassign(ctx context.Context, pullRequestId string, oldReviewerId string) (*entities.PullRequest, error) {
	pr, err := s.PrRepository.GetByID(ctx, pullRequestId)

	if err != nil {
		return nil, customerrors.NewDomainError(customerrors.NotFound, "PR not found")
	}

	if pr.Status == entities.StatusMerged {
		return nil, customerrors.NewDomainError(customerrors.PRMerged, "cannot reassign on merged PR")
	}
	author, err := s.UserRepository.GetById(ctx, pr.AuthorId)
	if err != nil {
		return nil, err
	}

	assigned := slices.Contains(pr.AssignedReviewers, oldReviewerId)

	if !assigned {
		return nil, customerrors.NewDomainError(customerrors.NotAssigned, "reviewer is not assigned to this PR")
	}

	newAssignedReviewers := pr.AssignedReviewers
	team, err := s.UserRepository.GetUserTeam(ctx, author.TeamName)

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

	for _, user := range team {
		if user != oldReviewerId && user != reviewerOne && user != reviewerTwo && user != pr.AuthorId {
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

	if len(newAssignedReviewers) == 0 {
		return nil, customerrors.NewDomainError(customerrors.NoCandidate, "no active replacement candidate in team")
	}

	err = s.PrRepository.Reassign(ctx, pullRequestId, newAssignedReviewers)
	pr, _ = s.PrRepository.GetByID(ctx, pullRequestId)
	return pr, err
}
