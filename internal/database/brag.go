// Package database provides database connection and migration management.
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

type sqliteBragRepository struct {
	db       *SQLiteDB
	userRepo repository.UserRepository
}

// NewBragRepository creates a new SQLite brag repository
func NewBragRepository(db *SQLiteDB, userRepo repository.UserRepository) repository.BragRepository {
	return &sqliteBragRepository{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *sqliteBragRepository) Select(ctx context.Context, id int64) (*domain.Brag, error) {
	dbBrag, err := r.db.Queries().GetBrag(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("brag not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get brag: %w", err)
	}

	return r.toDomainBrag(ctx, &dbBrag)
}

func (r *sqliteBragRepository) SelectAll(ctx context.Context, userID int64) ([]*domain.Brag, error) {
	dbBrags, err := r.db.Queries().ListBragsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brags: %w", err)
	}

	brags := make([]*domain.Brag, 0, len(dbBrags))
	for _, dbBrag := range dbBrags {
		brag, err := r.toDomainBrag(ctx, &dbBrag)
		if err != nil {
			return nil, err
		}
		brags = append(brags, brag)
	}

	return brags, nil
}

func (r *sqliteBragRepository) SelectByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error) {
	dbBrags, err := r.db.Queries().SearchBragsByTags(ctx, queries.SearchBragsByTagsParams{
		OwnerID:  userID,
		TagNames: tagNames,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search brags by tags: %w", err)
	}

	brags := make([]*domain.Brag, 0, len(dbBrags))
	for _, dbBrag := range dbBrags {
		brag, err := r.toDomainBrag(ctx, &dbBrag)
		if err != nil {
			return nil, err
		}
		brags = append(brags, brag)
	}

	return brags, nil
}

func (r *sqliteBragRepository) SelectByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error) {
	dbBrags, err := r.db.Queries().ListBragsByCategory(ctx, queries.ListBragsByCategoryParams{
		OwnerID:  userID,
		Category: int64(category),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list brags by category: %w", err)
	}

	brags := make([]*domain.Brag, 0, len(dbBrags))
	for _, dbBrag := range dbBrags {
		brag, err := r.toDomainBrag(ctx, &dbBrag)
		if err != nil {
			return nil, err
		}
		brags = append(brags, brag)
	}

	return brags, nil
}

func (r *sqliteBragRepository) Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	dbBrag, err := r.db.Queries().CreateBrag(ctx, queries.CreateBragParams{
		OwnerID:     brag.Owner.ID,
		Title:       brag.Title,
		Description: brag.Description,
		Category:    int64(brag.Category),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create brag: %w", err)
	}

	return r.toDomainBrag(ctx, &dbBrag)
}

func (r *sqliteBragRepository) Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	dbBrag, err := r.db.Queries().UpdateBrag(ctx, queries.UpdateBragParams{
		Title:       brag.Title,
		Description: brag.Description,
		Category:    int64(brag.Category),
		ID:          brag.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update brag: %w", err)
	}

	return r.toDomainBrag(ctx, &dbBrag)
}

func (r *sqliteBragRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.Queries().DeleteBrag(ctx, id); err != nil {
		return fmt.Errorf("failed to delete brag: %w", err)
	}
	return nil
}

func (r *sqliteBragRepository) toDomainBrag(ctx context.Context, dbBrag *queries.Brag) (*domain.Brag, error) {
	user, err := r.userRepo.Select(ctx, dbBrag.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("owner not found for brag %d: %w", dbBrag.ID, err)
		}
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	brag := &domain.Brag{
		ID:          dbBrag.ID,
		Owner:       *user,
		Title:       dbBrag.Title,
		Description: dbBrag.Description,
		Category:    domain.Category(dbBrag.Category),
	}

	if dbBrag.CreatedAt.Valid {
		brag.CreatedAt = dbBrag.CreatedAt.Time
	}

	if dbBrag.UpdatedAt.Valid {
		brag.UpdatedAt = dbBrag.UpdatedAt.Time
	}

	return brag, nil
}
