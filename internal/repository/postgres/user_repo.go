package postgres

import (
	"context"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, entity *entities.User) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO users (user_id, username, team_name, is_active) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) 
        DO UPDATE SET username = $2, team_name = $3, is_active = $4
    `, entity.UserId, entity.Username, entity.TeamName, entity.IsActive)
	return err
}

func (r *UserRepository) SetIsActive(ctx context.Context, userId string, isActive bool) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(ctx, `
        UPDATE users 
        SET is_active = $1 
        WHERE user_id = $2 
        RETURNING user_id, username, team_name, is_active
    `, isActive, userId).Scan(
		&user.UserId, &user.Username, &user.TeamName, &user.IsActive,
	)
	if err != nil {
		return nil, customerrors.NewDomainError(customerrors.NotFound, "user not found")
	}
	return &user, nil
}

func (r *UserRepository) GetById(ctx context.Context, userId string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(ctx, `
        SELECT user_id, username, team_name, is_active
        FROM users WHERE user_id = $1
    `, userId).Scan(
		&user.UserId, &user.Username, &user.TeamName, &user.IsActive,
	)
	if err != nil {
		return nil, customerrors.NewDomainError(customerrors.NotFound, "user not found")
	}
	return &user, nil
}

func (r *UserRepository) GetUserTeam(ctx context.Context, teamName string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
        SELECT user_id, username, team_name, is_active 
        FROM users 
        WHERE team_name = $1 AND is_active = true
    `, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user.UserId)
	}
	return users, nil
}
