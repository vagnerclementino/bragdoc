package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

type sqliteTagRepository struct {
	db *SQLiteDB
}

// NewTagRepository creates a new SQLite tag repository
func NewTagRepository(db *SQLiteDB) repository.TagRepository {
	return &sqliteTagRepository{db: db}
}

func (r *sqliteTagRepository) Select(ctx context.Context, id int64) (*domain.Tag, error) {
	dbTag, err := r.db.Queries().GetTag(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return toDomainTag(&dbTag), nil
}

func (r *sqliteTagRepository) SelectByName(ctx context.Context, ownerID int64, name string) (*domain.Tag, error) {
	dbTag, err := r.db.Queries().GetTagByName(ctx, queries.GetTagByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}

	return toDomainTag(&dbTag), nil
}

func (r *sqliteTagRepository) SelectAll(ctx context.Context, ownerID int64) ([]*domain.Tag, error) {
	dbTags, err := r.db.Queries().ListTagsByUser(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tags := make([]*domain.Tag, len(dbTags))
	for i, dbTag := range dbTags {
		tags[i] = toDomainTag(&dbTag)
	}

	return tags, nil
}

func (r *sqliteTagRepository) SelectByBrag(ctx context.Context, bragID int64) ([]*domain.Tag, error) {
	dbTags, err := r.db.Queries().ListTagsByBrag(ctx, bragID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags by brag: %w", err)
	}

	tags := make([]*domain.Tag, len(dbTags))
	for i, dbTag := range dbTags {
		tags[i] = toDomainTag(&dbTag)
	}

	return tags, nil
}

func (r *sqliteTagRepository) Insert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	dbTag, err := r.db.Queries().CreateTag(ctx, queries.CreateTagParams{
		Name:    tag.Name,
		OwnerID: tag.OwnerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return toDomainTag(&dbTag), nil
}

func (r *sqliteTagRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.Queries().DetachTagFromAllBrags(ctx, id); err != nil {
		return fmt.Errorf("failed to detach tag from brags: %w", err)
	}

	if err := r.db.Queries().DeleteTag(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (r *sqliteTagRepository) AttachToBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		if err := r.db.Queries().AttachTagToBrag(ctx, queries.AttachTagToBragParams{
			BragID: bragID,
			TagID:  tagID,
		}); err != nil {
			return fmt.Errorf("failed to attach tag %d to brag %d: %w", tagID, bragID, err)
		}
	}
	return nil
}

func (r *sqliteTagRepository) DetachFromBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		if err := r.db.Queries().DetachTagFromBrag(ctx, queries.DetachTagFromBragParams{
			BragID: bragID,
			TagID:  tagID,
		}); err != nil {
			return fmt.Errorf("failed to detach tag %d from brag %d: %w", tagID, bragID, err)
		}
	}
	return nil
}

func toDomainTag(dbTag *queries.Tag) *domain.Tag {
	tag := &domain.Tag{
		ID:      dbTag.ID,
		Name:    dbTag.Name,
		OwnerID: dbTag.OwnerID,
	}

	if dbTag.CreatedAt.Valid {
		tag.CreatedAt = dbTag.CreatedAt.Time
	}

	return tag
}
