package repository

import (
    "context"
    "github.com/vagnerclementino/bragdoc/internal/domain"
)

type CategoryRepository interface {
    Get(ctx context.Context, id int64) (*domain.Category, error)
    GetByName(ctx context.Context, name domain.CategoryName) (*domain.Category, error)
    List(ctx context.Context) ([]*domain.Category, error)
    Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
    Update(ctx context.Context, category *domain.Category) (*domain.Category, error)
    Delete(ctx context.Context, id int64) error
}
