package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// BragRepository defines the interface for brag data access
type BragRepository interface {
	Select(ctx context.Context, id int64) (*domain.Brag, error)
	SelectAll(ctx context.Context, userID int64) ([]*domain.Brag, error)
	SelectByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error)
	SelectByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error)
	Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error)
	Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error)
	Delete(ctx context.Context, id int64) error
}

type bragRepo struct {
	queries *queries.Queries
}

// NewBragRepository creates a new brag repository
func NewBragRepository(db *sql.DB) BragRepository {
	return &bragRepo{
		queries: queries.New(db),
	}
}

func (r *bragRepo) Select(ctx context.Context, id int64) (*domain.Brag, error) {
	dbBrag, err := r.queries.GetBrag(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("brag not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get brag: %w", err)
	}

	return r.toDomainBrag(&dbBrag), nil
}

func (r *bragRepo) SelectAll(ctx context.Context, userID int64) ([]*domain.Brag, error) {
	dbBrags, err := r.queries.ListBragsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brags: %w", err)
	}

	brags := make([]*domain.Brag, len(dbBrags))
	for i, dbBrag := range dbBrags {
		brags[i] = r.toDomainBrag(&dbBrag)
	}

	return brags, nil
}

func (r *bragRepo) SelectByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error) {
	dbBrags, err := r.queries.SearchBragsByTags(ctx, queries.SearchBragsByTagsParams{
		OwnerID:  userID,
		TagNames: tagNames,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search brags by tags: %w", err)
	}

	brags := make([]*domain.Brag, len(dbBrags))
	for i, dbBrag := range dbBrags {
		brags[i] = r.toDomainBrag(&dbBrag)
	}

	return brags, nil
}

func (r *bragRepo) SelectByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error) {
	dbBrags, err := r.queries.ListBragsByCategory(ctx, queries.ListBragsByCategoryParams{
		OwnerID:  userID,
		Category: int64(category),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list brags by category: %w", err)
	}

	brags := make([]*domain.Brag, len(dbBrags))
	for i, dbBrag := range dbBrags {
		brags[i] = r.toDomainBrag(&dbBrag)
	}

	return brags, nil
}

func (r *bragRepo) Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	dbBrag, err := r.queries.CreateBrag(ctx, queries.CreateBragParams{
		OwnerID:     brag.OwnerID,
		Title:       brag.Title,
		Description: brag.Description,
		Category:    int64(brag.Category),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create brag: %w", err)
	}

	return r.toDomainBrag(&dbBrag), nil
}

func (r *bragRepo) Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	dbBrag, err := r.queries.UpdateBrag(ctx, queries.UpdateBragParams{
		Title:       brag.Title,
		Description: brag.Description,
		Category:    int64(brag.Category),
		ID:          brag.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update brag: %w", err)
	}

	return r.toDomainBrag(&dbBrag), nil
}

func (r *bragRepo) Delete(ctx context.Context, id int64) error {
	if err := r.queries.DeleteBrag(ctx, id); err != nil {
		return fmt.Errorf("failed to delete brag: %w", err)
	}
	return nil
}

// toDomainBrag converts a database brag to a domain brag
func (r *bragRepo) toDomainBrag(dbBrag *queries.Brag) *domain.Brag {
	brag := &domain.Brag{
		ID:          dbBrag.ID,
		OwnerID:     dbBrag.OwnerID,
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

	return brag
}
