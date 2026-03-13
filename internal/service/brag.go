package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

type BragService struct {
	repo repository.BragRepository
}

func NewBragService(repo repository.BragRepository) *BragService {
	return &BragService{repo: repo}
}

func (s *BragService) validateBrag(brag *domain.Brag) error {
	if brag == nil {
		return errors.New("brag cannot be nil")
	}

	// Validate Owner
	if brag.Owner.ID == 0 {
		return errors.New("brag owner ID cannot be empty")
	}

	// Validate Category
	if err := brag.Category.Validate(); err != nil {
		return fmt.Errorf("invalid category: %w", err)
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

	return nil
}

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

func (s *BragService) GetByID(ctx context.Context, id int64) (*domain.Brag, error) {
	return s.repo.Select(ctx, id)
}

func (s *BragService) List(ctx context.Context, userID int64) ([]*domain.Brag, error) {
	return s.repo.SelectAll(ctx, userID)
}

func (s *BragService) SearchByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error) {
	if len(tagNames) == 0 {
		return nil, fmt.Errorf("at least one tag name is required")
	}
	return s.repo.SelectByTags(ctx, userID, tagNames)
}

func (s *BragService) SearchByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error) {
	return s.repo.SelectByCategory(ctx, userID, category)
}

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

func (s *BragService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
