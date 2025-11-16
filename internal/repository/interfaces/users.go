package interfaces

import (
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository"
)

type UserRepository interface {
	repository.Repository[entities.User]

	SetIsActive(userId string, isActive bool) (*entities.User, error)
	GetUserTeam(userId string) (*[]string, error)
	GetById(userId string) (*entities.User, error)
}
