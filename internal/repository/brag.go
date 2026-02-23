package repository

import (
    "context"
    "github.com/vagnerclementino/bragdoc/internal/domain"
)

type BragRepository interface {
    Select(ctx context.Context, id int64) (*domain.Brag, error)
    SelectAll(ctx context.Context, userID int64) ([]*domain.Brag, error)
    SelectByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error)
    SelectByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error)
    Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error)
    Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error)
    Delete(ctx context.Context, id int64) error
}
