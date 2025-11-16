package inmemory

import (
	"context"
	"fmt"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/repository/interfaces"
)

type UserRepository struct {
	users []entities.User
}

func NewUserRepository() interfaces.UserRepository {
	repo := &UserRepository{}
	user_one := entities.NewUser("u1", "username", "team", true)
	user_two := entities.NewUser("u2", "username", "team", true)
	user_three := entities.NewUser("u3", "username", "team", true)
	user_four := entities.NewUser("u4", "username", "team", true)
	repo.users = []entities.User{*user_one, *user_two, *user_three, *user_four}
	return repo
}

func (r *UserRepository) Create(ctx context.Context, entity *entities.User) error {
	r.users = append(r.users, *entity)
	return nil
}

func (r *UserRepository) SetIsActive(ctx context.Context, userId string, isActive bool) (*entities.User, error) {
	user, err := r.GetById(ctx, userId)
	user.IsActive = isActive

	return user, err
}

func (r *UserRepository) GetById(ctx context.Context, userId string) (*entities.User, error) {
	for i := range r.users {
		if r.users[i].UserId == userId {
			return &r.users[i], nil
		}
	}

	return nil, fmt.Errorf("no user with Id '%s'", userId)
}

func (r *UserRepository) GetUserTeam(ctx context.Context, teamName string) ([]string, error) {
	var teamMembers []string
	for i := range r.users {
		if r.users[i].TeamName == teamName && r.users[i].IsActive {
			teamMembers = append(teamMembers, r.users[i].UserId)
		}
	}

	return teamMembers, nil
}
