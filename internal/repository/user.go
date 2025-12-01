package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
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

type userRepo struct {
	queries *queries.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		queries: queries.New(db),
	}
}

func (r *userRepo) Select(ctx context.Context, id int64) (*domain.User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

func (r *userRepo) SelectByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

func (r *userRepo) SelectAll(ctx context.Context) ([]*domain.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(&dbUser)
	}

	return users, nil
}

func (r *userRepo) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	jobTitle := sql.NullString{String: user.JobTitle, Valid: user.JobTitle != ""}
	company := sql.NullString{String: user.Company, Valid: user.Company != ""}

	dbUser, err := r.queries.CreateUser(ctx, queries.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		JobTitle: jobTitle,
		Company:  company,
		Locale:   user.Locale.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	jobTitle := sql.NullString{String: user.JobTitle, Valid: user.JobTitle != ""}
	company := sql.NullString{String: user.Company, Valid: user.Company != ""}

	dbUser, err := r.queries.UpdateUser(ctx, queries.UpdateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		JobTitle: jobTitle,
		Company:  company,
		Locale:   user.Locale.String(),
		ID:       user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.toDomainUser(&dbUser), nil
}

func (r *userRepo) Delete(ctx context.Context, id int64) error {
	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// toDomainUser converts a database user to a domain user
func (r *userRepo) toDomainUser(dbUser *queries.User) *domain.User {
	locale, _ := domain.ParseLocale(dbUser.Locale)
	
	user := &domain.User{
		ID:     dbUser.ID,
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Locale: locale,
	}

	if dbUser.JobTitle.Valid {
		user.JobTitle = dbUser.JobTitle.String
	}

	if dbUser.Company.Valid {
		user.Company = dbUser.Company.String
	}

	if dbUser.CreatedAt.Valid {
		user.CreatedAt = dbUser.CreatedAt.Time
	}

	if dbUser.UpdatedAt.Valid {
		user.UpdatedAt = dbUser.UpdatedAt.Time
	}

	return user
}
