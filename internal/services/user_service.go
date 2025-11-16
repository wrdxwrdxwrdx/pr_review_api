package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type UserService struct {
	userRepository interfaces.UserRepository
	prRepository   interfaces.PrRepository
}

func NewUserService(userRepository interfaces.UserRepository, prRepository interfaces.PrRepository) UserService {
	return UserService{userRepository: userRepository, prRepository: prRepository}
}

func (s *UserService) Create(ctx context.Context, username string, teamName string) (*entities.User, error) {
	user := entities.NewUser("u1", username, teamName, true)
	s.userRepository.Create(ctx, user)
	return user, nil
}

func (s *UserService) SetIsActive(ctx context.Context, userId string, isActive bool) (*entities.User, error) {
	user, err := s.userRepository.SetIsActive(ctx, userId, isActive)
	return user, err
}

func (s *UserService) GetReview(ctx context.Context, userId string) (*entities.Reviews, error) {
	prs, err := s.prRepository.GetReview(ctx, userId)

	if err != nil {
		return nil, err
	}

	return entities.NewReviews(userId, prs), nil
}
