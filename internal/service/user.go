package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

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

// emailRegex is a simple regex for basic email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateUser performs comprehensive business validation
func (s *UserService) validateUser(user *domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// Normalize and validate
	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return errors.New("user name cannot be empty")
	}

	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" {
		return errors.New("user email cannot be empty")
	}

	// Business validations
	if len(user.Name) < 2 {
		return fmt.Errorf("user name must be at least 2 characters, got %d", len(user.Name))
	}

	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("invalid email format: %s", user.Email)
	}

	// Default locale to en-US if empty
	if user.Locale == "" {
		user.Locale = domain.LocaleEnglishUS
	}

	if !user.Locale.IsValid() {
		return fmt.Errorf("invalid locale: %s (supported: %s)", user.Locale, domain.SupportedLocalesString())
	}

	return nil
}

// Create creates a new user with validation
func (s *UserService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validate user data
	if err := s.validateUser(user); err != nil {
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
	if err := s.validateUser(user); err != nil {
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
