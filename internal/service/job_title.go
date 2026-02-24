package service

import (
    "context"
    "time"
    "github.com/vagnerclementino/bragdoc/internal/domain"
    "github.com/vagnerclementino/bragdoc/internal/repository"
)

type JobTitleService struct {
    repo repository.JobTitleRepository
}

func NewJobTitleService(repo repository.JobTitleRepository) *JobTitleService {
    return &JobTitleService{repo: repo}
}

func (s *JobTitleService) GetOrCreate(ctx context.Context, userID int64, title string, company string) (*domain.JobTitle, error) {
    // Try to get existing job title by name
    jobTitle, err := s.repo.GetByName(ctx, userID, title)
    if err == nil {
        return jobTitle, nil
    }

    // Get current active job title to close it
    activeJobTitle, err := s.repo.GetActive(ctx, userID)
    if err == nil && activeJobTitle != nil {
        // Close the current active job title
        now := time.Now()
        activeJobTitle.EndDate = &now
        _, _ = s.repo.Update(ctx, activeJobTitle)
    }

    // Create new active job title
    user := domain.User{ID: userID}
    newJobTitle := &domain.JobTitle{
        User:    user,
        Title:   title,
        Company: company,
    }

    return s.repo.Create(ctx, newJobTitle)
}

func (s *JobTitleService) GetActive(ctx context.Context, userID int64) (*domain.JobTitle, error) {
    return s.repo.GetActive(ctx, userID)
}

func (s *JobTitleService) GetByID(ctx context.Context, id int64) (*domain.JobTitle, error) {
    return s.repo.Get(ctx, id)
}
