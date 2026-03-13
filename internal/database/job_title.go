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

type sqliteJobTitleRepository struct {
	db       *SQLiteDB
	userRepo repository.UserRepository
}

func NewJobTitleRepository(db *SQLiteDB, userRepo repository.UserRepository) repository.JobTitleRepository {
	return &sqliteJobTitleRepository{db: db, userRepo: userRepo}
}

func (r *sqliteJobTitleRepository) Get(ctx context.Context, id int64) (*domain.JobTitle, error) {
	dbJobTitle, err := r.db.Queries().GetJobTitle(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("job title not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get job title: %w", err)
	}
	return r.toDomainJobTitle(ctx, &dbJobTitle)
}

func (r *sqliteJobTitleRepository) GetActive(ctx context.Context, userID int64) (*domain.JobTitle, error) {
	dbJobTitle, err := r.db.Queries().GetActiveJobTitle(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active job title found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to get active job title: %w", err)
	}
	return r.toDomainJobTitle(ctx, &dbJobTitle)
}

func (r *sqliteJobTitleRepository) GetByName(ctx context.Context, userID int64, title string) (*domain.JobTitle, error) {
	dbJobTitle, err := r.db.Queries().GetJobTitleByName(ctx, queries.GetJobTitleByNameParams{
		UserID: userID,
		Title:  title,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("job title '%s' not found for user %d", title, userID)
		}
		return nil, fmt.Errorf("failed to get job title by name: %w", err)
	}
	return r.toDomainJobTitle(ctx, &dbJobTitle)
}

func (r *sqliteJobTitleRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.JobTitle, error) {
	dbJobTitles, err := r.db.Queries().ListJobTitlesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list job titles: %w", err)
	}

	jobTitles := make([]*domain.JobTitle, 0, len(dbJobTitles))
	for _, dbJobTitle := range dbJobTitles {
		jobTitle, err := r.toDomainJobTitle(ctx, &dbJobTitle)
		if err != nil {
			return nil, err
		}
		jobTitles = append(jobTitles, jobTitle)
	}
	return jobTitles, nil
}

func (r *sqliteJobTitleRepository) Create(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error) {
	if err := r.validateJobTitle(jobTitle); err != nil {
		return nil, fmt.Errorf("invalid job title: %w", err)
	}

	var startDate, endDate sql.NullTime
	if jobTitle.StartDate != nil {
		startDate = sql.NullTime{Time: *jobTitle.StartDate, Valid: true}
	}
	if jobTitle.EndDate != nil {
		endDate = sql.NullTime{Time: *jobTitle.EndDate, Valid: true}
	}

	dbJobTitle, err := r.db.Queries().CreateJobTitle(ctx, queries.CreateJobTitleParams{
		UserID:    jobTitle.User.ID,
		Title:     jobTitle.Title,
		Company:   sql.NullString{String: jobTitle.Company, Valid: jobTitle.Company != ""},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create job title: %w", err)
	}
	return r.toDomainJobTitle(ctx, &dbJobTitle)
}

func (r *sqliteJobTitleRepository) Update(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error) {
	if err := r.validateJobTitle(jobTitle); err != nil {
		return nil, fmt.Errorf("invalid job title: %w", err)
	}

	var startDate, endDate sql.NullTime
	if jobTitle.StartDate != nil {
		startDate = sql.NullTime{Time: *jobTitle.StartDate, Valid: true}
	}
	if jobTitle.EndDate != nil {
		endDate = sql.NullTime{Time: *jobTitle.EndDate, Valid: true}
	}

	dbJobTitle, err := r.db.Queries().UpdateJobTitle(ctx, queries.UpdateJobTitleParams{
		Title:     jobTitle.Title,
		Company:   sql.NullString{String: jobTitle.Company, Valid: jobTitle.Company != ""},
		StartDate: startDate,
		EndDate:   endDate,
		ID:        jobTitle.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update job title: %w", err)
	}
	return r.toDomainJobTitle(ctx, &dbJobTitle)
}

func (r *sqliteJobTitleRepository) Delete(ctx context.Context, id int64) error {
	// Check if job title is in use
	count, err := r.db.Queries().CountBragsByJobTitle(ctx, sql.NullInt64{Int64: id, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to check job title usage: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete job title %d: %d brags use it", id, count)
	}

	if err := r.db.Queries().DeleteJobTitle(ctx, id); err != nil {
		return fmt.Errorf("failed to delete job title: %w", err)
	}
	return nil
}

func (r *sqliteJobTitleRepository) validateJobTitle(jobTitle *domain.JobTitle) error {
	if jobTitle == nil {
		return errors.New("job title cannot be nil")
	}

	if jobTitle.User.ID == 0 {
		return errors.New("job title user ID cannot be empty")
	}

	if jobTitle.Title == "" {
		return errors.New("job title cannot be empty")
	}

	if jobTitle.Company == "" {
		return errors.New("company cannot be empty")
	}

	// Validate date range if both dates are set
	if jobTitle.StartDate != nil && jobTitle.EndDate != nil {
		if jobTitle.EndDate.Before(*jobTitle.StartDate) {
			return errors.New("end date cannot be before start date")
		}
	}

	return nil
}

func (r *sqliteJobTitleRepository) toDomainJobTitle(ctx context.Context, dbJobTitle *queries.JobTitle) (*domain.JobTitle, error) {
	user, err := r.userRepo.Select(ctx, dbJobTitle.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found for job title %d: %w", dbJobTitle.ID, err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	jobTitle := &domain.JobTitle{
		ID:    dbJobTitle.ID,
		User:  *user,
		Title: dbJobTitle.Title,
	}

	if dbJobTitle.Company.Valid {
		jobTitle.Company = dbJobTitle.Company.String
	}

	if dbJobTitle.StartDate.Valid {
		startDate := dbJobTitle.StartDate.Time
		jobTitle.StartDate = &startDate
	}

	if dbJobTitle.EndDate.Valid {
		endDate := dbJobTitle.EndDate.Time
		jobTitle.EndDate = &endDate
	}

	if dbJobTitle.CreatedAt.Valid {
		jobTitle.CreatedAt = dbJobTitle.CreatedAt.Time
	}
	if dbJobTitle.UpdatedAt.Valid {
		jobTitle.UpdatedAt = dbJobTitle.UpdatedAt.Time
	}

	return jobTitle, nil
}
