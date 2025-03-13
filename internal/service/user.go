package service

import (
	"context"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) InitializeUser(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		ID:    generateID(),
		Name:  name,
		Email: email,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context) (*domain.User, error) {
	return s.repo.Find(ctx)
}
