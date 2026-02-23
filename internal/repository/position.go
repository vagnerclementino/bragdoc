package repository

import (
    "context"
    "github.com/vagnerclementino/bragdoc/internal/domain"
)

type PositionRepository interface {
    Get(ctx context.Context, id int64) (*domain.Position, error)
    ListByUser(ctx context.Context, userID int64) ([]*domain.Position, error)
    Create(ctx context.Context, position *domain.Position) (*domain.Position, error)
    Update(ctx context.Context, position *domain.Position) (*domain.Position, error)
    Delete(ctx context.Context, id int64) error
}
