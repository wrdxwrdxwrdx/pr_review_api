package services

import (
	"context"
	"fmt"
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

func (s *UserService) Create(ctx context.Context, userId string, username string, teamName string, isActive bool) (*entities.User, error) {
	user := entities.NewUser(userId, username, teamName, isActive)
	fmt.Println(user)
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

func (s *UserService) GetById(ctx context.Context, userId string) (*entities.User, error) {
	user, err := s.userRepository.GetById(ctx, userId)
	return user, err
}

func (s *UserService) Exist(ctx context.Context, userId string) (bool, error) {
	user, _ := s.userRepository.GetById(ctx, userId)
	fmt.Println(user)
	if user != nil {
		return true, nil
	}
	return false, nil
}
