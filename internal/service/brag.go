package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

// BragService provides business logic for brag management
type BragService struct {
	repo repository.BragRepository
}

// NewBragService creates a new brag service
func NewBragService(repo repository.BragRepository) *BragService {
	return &BragService{repo: repo}
}

// validateBrag performs comprehensive business validation
func (s *BragService) validateBrag(brag *domain.Brag) error {
	if brag == nil {
		return errors.New("brag cannot be nil")
	}

	// Structural validations
	title := strings.TrimSpace(brag.Title)
	if title == "" {
		return errors.New("brag title cannot be empty")
	}

	description := strings.TrimSpace(brag.Description)
	if description == "" {
		return errors.New("brag description cannot be empty")
	}

	// Business validations
	if len(title) < 5 {
		return fmt.Errorf("brag title must be at least 5 characters, got %d", len(title))
	}

	if len(description) < 20 {
		return fmt.Errorf("brag description must be at least 20 characters, got %d", len(description))
	}

	if brag.Category < domain.CategoryProject || brag.Category > domain.CategoryInnovation {
		return fmt.Errorf("invalid brag category: %d", brag.Category)
	}

	return nil
}

// Create creates a new brag with validation
func (s *BragService) Create(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	if err := s.validateBrag(brag); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	created, err := s.repo.Insert(ctx, brag)
	if err != nil {
		return nil, fmt.Errorf("failed to create brag: %w", err)
	}

	return created, nil
}

// GetByID retrieves a brag by ID
func (s *BragService) GetByID(ctx context.Context, id int64) (*domain.Brag, error) {
	return s.repo.Select(ctx, id)
}

// List retrieves all brags for a user
func (s *BragService) List(ctx context.Context, userID int64) ([]*domain.Brag, error) {
	return s.repo.SelectAll(ctx, userID)
}

// SearchByTags retrieves brags filtered by tags
func (s *BragService) SearchByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error) {
	if len(tagNames) == 0 {
		return nil, fmt.Errorf("at least one tag name is required")
	}
	return s.repo.SelectByTags(ctx, userID, tagNames)
}

// SearchByCategory retrieves brags filtered by category
func (s *BragService) SearchByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error) {
	return s.repo.SelectByCategory(ctx, userID, category)
}

// Update updates an existing brag
func (s *BragService) Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	if err := s.validateBrag(brag); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updated, err := s.repo.Update(ctx, brag)
	if err != nil {
		return nil, fmt.Errorf("failed to update brag: %w", err)
	}

	return updated, nil
}

// Delete deletes a brag by ID
func (s *BragService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
