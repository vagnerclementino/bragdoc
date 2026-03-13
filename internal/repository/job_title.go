package repository

import (
	"context"

	"github.com/vagnerclementino/bragdoc/internal/domain"
)

type JobTitleRepository interface {
	Get(ctx context.Context, id int64) (*domain.JobTitle, error)
	GetActive(ctx context.Context, userID int64) (*domain.JobTitle, error)
	GetByName(ctx context.Context, userID int64, title string) (*domain.JobTitle, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.JobTitle, error)
	Create(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error)
	Update(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error)
	Delete(ctx context.Context, id int64) error
}
