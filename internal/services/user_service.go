package services

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(repository interfaces.UserRepository) UserService {
	return UserService{userRepository: repository}
}

func (s UserService) Create(ctx context.Context, username string, teamName string) (*entities.User, error) {
	user := entities.NewUser(uuid.New(), username, teamName, true)
	s.userRepository.Create(ctx, user)
	return user, nil
}
