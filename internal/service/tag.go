package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

// TagService provides business logic for tag management
type TagService struct {
	repo repository.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// validateTag performs comprehensive business validation
func (s *TagService) validateTag(tag *domain.Tag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}

	// Structural validations
	name := strings.TrimSpace(tag.Name)
	if name == "" {
		return errors.New("tag name cannot be empty")
	}

	// Business validations
	if len(name) < 2 {
		return fmt.Errorf("tag name must be at least 2 characters, got %d", len(name))
	}

	if len(name) > 20 {
		return fmt.Errorf("tag name cannot exceed 20 characters, got %d", len(name))
	}

	return nil
}

// Create creates a new tag with validation
func (s *TagService) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	if err := s.validateTag(tag); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if tag already exists for this user
	existing, err := s.repo.SelectByName(ctx, tag.OwnerID, tag.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("tag '%s' already exists for this user", tag.Name)
	}

	created, err := s.repo.Insert(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return created, nil
}

// GetByID retrieves a tag by ID
func (s *TagService) GetByID(ctx context.Context, id int64) (*domain.Tag, error) {
	return s.repo.Select(ctx, id)
}

// GetByName retrieves a tag by name for a specific user
func (s *TagService) GetByName(ctx context.Context, ownerID int64, name string) (*domain.Tag, error) {
	return s.repo.SelectByName(ctx, ownerID, name)
}

// ListByUser retrieves all tags for a user
func (s *TagService) ListByUser(ctx context.Context, ownerID int64) ([]*domain.Tag, error) {
	return s.repo.SelectAll(ctx, ownerID)
}

// ListByBrag retrieves all tags associated with a brag
func (s *TagService) ListByBrag(ctx context.Context, bragID int64) ([]*domain.Tag, error) {
	return s.repo.SelectByBrag(ctx, bragID)
}

// AttachToBrag attaches tags to a brag
func (s *TagService) AttachToBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return fmt.Errorf("at least one tag ID is required")
	}
	return s.repo.AttachToBrag(ctx, bragID, tagIDs)
}

// DetachFromBrag detaches tags from a brag
func (s *TagService) DetachFromBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return fmt.Errorf("at least one tag ID is required")
	}
	return s.repo.DetachFromBrag(ctx, bragID, tagIDs)
}

// Delete deletes a tag by ID
func (s *TagService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetOrCreate gets an existing tag or creates a new one if it doesn't exist
func (s *TagService) GetOrCreate(ctx context.Context, ownerID int64, name string) (*domain.Tag, error) {
	// Try to get existing tag
	existing, err := s.repo.SelectByName(ctx, ownerID, name)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new tag
	tag := &domain.Tag{
		Name:    name,
		OwnerID: ownerID,
	}

	return s.Create(ctx, tag)
}
