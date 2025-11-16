package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Connect(connString string) error {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	db.Pool = pool
	log.Println("Successfully connected to database")
	return nil
}

func (db *DB) InitSchema(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS teams (
            team_name VARCHAR(255) PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,

		`CREATE TABLE IF NOT EXISTS users (
            user_id VARCHAR(255) PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            team_name VARCHAR(255) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
            is_active BOOLEAN NOT NULL DEFAULT true,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,

		`CREATE TABLE IF NOT EXISTS pull_requests (
            pull_request_id VARCHAR(255) PRIMARY KEY,
            pull_request_name VARCHAR(255) NOT NULL,
            author_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
            status VARCHAR(50) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
            assigned_reviewers TEXT[],
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            merged_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,

		`CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name)`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_created_at ON pull_requests(created_at)`,

		`CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $$ language 'plpgsql'`,

		`DROP TRIGGER IF EXISTS update_users_updated_at ON users`,
		`CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
         FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_pull_requests_updated_at ON pull_requests`,
		`CREATE TRIGGER update_pull_requests_updated_at BEFORE UPDATE ON pull_requests
         FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,
	}

	for _, query := range queries {
		_, err := db.Pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	log.Println("Database schema initialized successfully")
	return nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Database connection closed")
	}
}
