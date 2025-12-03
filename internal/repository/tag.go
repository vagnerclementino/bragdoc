package repository

import (
	"context"

	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	Select(ctx context.Context, id int64) (*domain.Tag, error)
	SelectByName(ctx context.Context, ownerID int64, name string) (*domain.Tag, error)
	SelectAll(ctx context.Context, ownerID int64) ([]*domain.Tag, error)
	SelectByBrag(ctx context.Context, bragID int64) ([]*domain.Tag, error)
	Insert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	Delete(ctx context.Context, id int64) error
	AttachToBrag(ctx context.Context, bragID int64, tagIDs []int64) error
	DetachFromBrag(ctx context.Context, bragID int64, tagIDs []int64) error
}
