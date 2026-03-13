package repository

import (
	"context"

	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Select(ctx context.Context, id int64) (*domain.User, error)
	SelectByEmail(ctx context.Context, email string) (*domain.User, error)
	SelectAll(ctx context.Context) ([]*domain.User, error)
	Insert(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}
