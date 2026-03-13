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

type sqliteCategoryRepository struct {
	db *SQLiteDB
}

func NewCategoryRepository(db *SQLiteDB) repository.CategoryRepository {
	return &sqliteCategoryRepository{db: db}
}

func (r *sqliteCategoryRepository) Get(ctx context.Context, id int64) (*domain.Category, error) {
	dbCategory, err := r.db.Queries().GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return r.toDomainCategory(&dbCategory), nil
}

func (r *sqliteCategoryRepository) GetByName(ctx context.Context, name domain.CategoryName) (*domain.Category, error) {
	dbCategory, err := r.db.Queries().GetCategoryByName(ctx, string(name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}
	return r.toDomainCategory(&dbCategory), nil
}

func (r *sqliteCategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	dbCategories, err := r.db.Queries().ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	categories := make([]*domain.Category, 0, len(dbCategories))
	for _, dbCategory := range dbCategories {
		categories = append(categories, r.toDomainCategory(&dbCategory))
	}
	return categories, nil
}

func (r *sqliteCategoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	dbCategory, err := r.db.Queries().CreateCategory(ctx, queries.CreateCategoryParams{
		Name:        string(category.Name),
		Description: sql.NullString{String: category.Description, Valid: category.Description != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return r.toDomainCategory(&dbCategory), nil
}

func (r *sqliteCategoryRepository) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	dbCategory, err := r.db.Queries().UpdateCategory(ctx, queries.UpdateCategoryParams{
		Name:        string(category.Name),
		Description: sql.NullString{String: category.Description, Valid: category.Description != ""},
		ID:          category.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	return r.toDomainCategory(&dbCategory), nil
}

func (r *sqliteCategoryRepository) Delete(ctx context.Context, id int64) error {
	// Check if category is in use
	count, err := r.db.Queries().CountBragsByCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category %d: %d brags use it", id, count)
	}

	if err := r.db.Queries().DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (r *sqliteCategoryRepository) toDomainCategory(dbCategory *queries.Category) *domain.Category {
	category := &domain.Category{
		ID:   dbCategory.ID,
		Name: domain.CategoryName(dbCategory.Name),
	}

	if dbCategory.Description.Valid {
		category.Description = dbCategory.Description.String
	}

	if dbCategory.CreatedAt.Valid {
		category.CreatedAt = dbCategory.CreatedAt.Time
	}
	if dbCategory.UpdatedAt.Valid {
		category.UpdatedAt = dbCategory.UpdatedAt.Time
	}

	return category
}
