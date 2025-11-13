package repository

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	// GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	// Update(ctx context.Context, entity *T) error
	// Delete(ctx context.Context, id uuid.UUID) error
	// List(ctx context.Context, limit, offset int) ([]*T, error)
}
