package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(repository interfaces.UserRepository) UserService {
	return UserService{userRepository: repository}
}

func (s *UserService) Create(ctx context.Context, username string, teamName string) (*entities.User, error) {
	user := entities.NewUser("u1", username, teamName, true)
	s.userRepository.Create(ctx, user)
	return user, nil
}

func (s *UserService) SetIsActive(ctx context.Context, userId string, isActive bool) (*entities.User, error) {
	user, err := s.userRepository.SetIsActive(userId, isActive)
	return user, err
}
