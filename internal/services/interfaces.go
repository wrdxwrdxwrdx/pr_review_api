package services

import "context"

type Service[T any] interface {
	Create(ctx context.Context, username string, team_name string, is_active bool) (*T, error)
}
