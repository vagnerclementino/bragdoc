package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
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

type tagRepo struct {
	queries *queries.Queries
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepo{
		queries: queries.New(db),
	}
}

func (r *tagRepo) Select(ctx context.Context, id int64) (*domain.Tag, error) {
	dbTag, err := r.queries.GetTag(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return r.toDomainTag(&dbTag), nil
}

func (r *tagRepo) SelectByName(ctx context.Context, ownerID int64, name string) (*domain.Tag, error) {
	dbTag, err := r.queries.GetTagByName(ctx, queries.GetTagByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}

	return r.toDomainTag(&dbTag), nil
}

func (r *tagRepo) SelectAll(ctx context.Context, ownerID int64) ([]*domain.Tag, error) {
	dbTags, err := r.queries.ListTagsByUser(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tags := make([]*domain.Tag, len(dbTags))
	for i, dbTag := range dbTags {
		tags[i] = r.toDomainTag(&dbTag)
	}

	return tags, nil
}

func (r *tagRepo) SelectByBrag(ctx context.Context, bragID int64) ([]*domain.Tag, error) {
	dbTags, err := r.queries.ListTagsByBrag(ctx, bragID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags by brag: %w", err)
	}

	tags := make([]*domain.Tag, len(dbTags))
	for i, dbTag := range dbTags {
		tags[i] = r.toDomainTag(&dbTag)
	}

	return tags, nil
}

func (r *tagRepo) Insert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	dbTag, err := r.queries.CreateTag(ctx, queries.CreateTagParams{
		Name:    tag.Name,
		OwnerID: tag.OwnerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return r.toDomainTag(&dbTag), nil
}

func (r *tagRepo) Delete(ctx context.Context, id int64) error {
	// First, detach tag from all brags
	if err := r.queries.DetachTagFromAllBrags(ctx, id); err != nil {
		return fmt.Errorf("failed to detach tag from brags: %w", err)
	}

	// Then delete the tag
	if err := r.queries.DeleteTag(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (r *tagRepo) AttachToBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		if err := r.queries.AttachTagToBrag(ctx, queries.AttachTagToBragParams{
			BragID: bragID,
			TagID:  tagID,
		}); err != nil {
			return fmt.Errorf("failed to attach tag %d to brag %d: %w", tagID, bragID, err)
		}
	}
	return nil
}

func (r *tagRepo) DetachFromBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		if err := r.queries.DetachTagFromBrag(ctx, queries.DetachTagFromBragParams{
			BragID: bragID,
			TagID:  tagID,
		}); err != nil {
			return fmt.Errorf("failed to detach tag %d from brag %d: %w", tagID, bragID, err)
		}
	}
	return nil
}

// toDomainTag converts a database tag to a domain tag
func (r *tagRepo) toDomainTag(dbTag *queries.Tag) *domain.Tag {
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
