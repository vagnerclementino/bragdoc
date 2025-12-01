package service

import (
	"context"
	"fmt"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

// UserService provides business logic for user management
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user with validation
func (s *UserService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validate user data
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	existing, err := s.repo.SelectByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Create user
	created, err := s.repo.Insert(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return created, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.Select(ctx, id)
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.SelectByEmail(ctx, email)
}

// List retrieves all users
func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.SelectAll(ctx)
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updated, nil
}

// Delete deletes a user by ID
func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
