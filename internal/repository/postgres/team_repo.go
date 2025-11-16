package postgres

import (
	"context"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"
	"pr_review_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db             *pgxpool.Pool
	userRepository UserRepository
}

func NewTeamRepository(db *pgxpool.Pool, userRepository UserRepository) *TeamRepository {
	return &TeamRepository{db: db, userRepository: userRepository}
}

func (r *TeamRepository) Create(ctx context.Context, entity *entities.Team) error {
	_, err := r.db.Exec(ctx, "INSERT INTO teams (team_name) VALUES ($1)", entity.TeamName)
	if err != nil {
		if utils.IsUnique(err) {
			return customerrors.NewDomainError(customerrors.TeamExists, "team already exists")
		}
		return err
	}

	for _, member := range entity.Members {
		user := entities.NewUser(member.UserId, member.Username, entity.TeamName, member.IsActive)
		err := r.userRepository.Create(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil
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
		return nil, customerrors.NewDomainError(customerrors.NotFound, "team not found")
	}

	return &team, nil
}
