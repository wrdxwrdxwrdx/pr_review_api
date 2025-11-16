package postgres

import (
	"context"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/domain/errors"
	"pr_review_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, entity *entities.Team) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO teams (team_name) VALUES ($1)", entity.TeamName)
	if err != nil {
		if utils.IsUnique(err) {
			return errors.NewDomainError(errors.TeamExists, "team already exists")
		}
		return err
	}

	for _, member := range entity.Members {
		_, err = tx.Exec(ctx, `
            INSERT INTO users (user_id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) 
            DO UPDATE SET username = $2, team_name = $3, is_active = $4
        `, member.UserId, member.Username, entity.TeamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) GetByName(ctx context.Context, teamNameQuery string) (*entities.Team, error) {
	var team entities.Team
	team.TeamName = teamNameQuery

	rows, err := r.db.Query(ctx, `
        SELECT user_id, username, is_active 
        FROM users 
        WHERE team_name = $1
    `, teamNameQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member entities.TeamMember
		err := rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	if len(team.Members) == 0 {
		return nil, errors.NewDomainError(errors.NotFound, "team not found")
	}

	return &team, nil
}
