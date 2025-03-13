package service

import (
	"context"
	"time"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

type BragService struct {
	repo repository.BragRepository
}

func NewBragService(repo repository.BragRepository) *BragService {
	return &BragService{repo: repo}
}

func (s *BragService) CreateBrag(ctx context.Context, description string, details *string) (*domain.Brag, error) {
	brag := &domain.Brag{
		ID:          generateID(),
		Description: description,
		Details:     details,
		CreatedAt:   time.Now().Unix(),
	}

	if err := brag.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, brag); err != nil {
		return nil, err
	}

	return brag, nil
}

func (s *BragService) ListBrags(ctx context.Context) ([]domain.Brag, error) {
	return s.repo.FindAll(ctx)
}

func generateID() string {
	// Implement your ID generation logic
	return "generated-id"
}
