package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

type sqliteUserRepository struct {
	db *SQLiteDB
}

// NewUserRepository creates a new SQLite user repository
func NewUserRepository(db *SQLiteDB) repository.UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Select(ctx context.Context, id int64) (*domain.User, error) {
	dbUser, err := r.db.Queries().GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %d: %w", id, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return toDomainUser(&dbUser), nil
}

func (r *sqliteUserRepository) SelectByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.db.Queries().GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s: %w", email, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return toDomainUser(&dbUser), nil
}

func (r *sqliteUserRepository) SelectAll(ctx context.Context) ([]*domain.User, error) {
	dbUsers, err := r.db.Queries().ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = toDomainUser(&dbUser)
	}

	return users, nil
}

func (r *sqliteUserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	jobTitle := sql.NullString{String: user.JobTitle, Valid: user.JobTitle != ""}
	company := sql.NullString{String: user.Company, Valid: user.Company != ""}

	dbUser, err := r.db.Queries().CreateUser(ctx, queries.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		JobTitle: jobTitle,
		Company:  company,
		Locale:   user.Locale.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toDomainUser(&dbUser), nil
}

func (r *sqliteUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	jobTitle := sql.NullString{String: user.JobTitle, Valid: user.JobTitle != ""}
	company := sql.NullString{String: user.Company, Valid: user.Company != ""}

	dbUser, err := r.db.Queries().UpdateUser(ctx, queries.UpdateUserParams{
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

	return toDomainUser(&dbUser), nil
}

func (r *sqliteUserRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.Queries().DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func toDomainUser(dbUser *queries.User) *domain.User {
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
