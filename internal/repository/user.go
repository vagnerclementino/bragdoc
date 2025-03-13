package repository

import (
	"context"

	"github.com/chaisql/chai"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Find(ctx context.Context) (*domain.User, error)
}

type userRepo struct {
	db *chai.DB
}

func NewUserRepository(db *chai.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.Exec(`
		INSERT INTO users (id, name, email) 
		VALUES (?, ?, ?)`,
		user.ID, user.Name, user.Email)
}

func (r *userRepo) Find(ctx context.Context) (*domain.User, error) {
	var user domain.User

	res, err := r.db.Query("SELECT * FROM users LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	err = res.Iterate(func(r *chai.Row) error {
		return r.Scan(&user.ID, &user.Name, &user.Email)
	})

	if err == nil && user.ID == "" {
		return nil, nil
	}

	return &user, err
}
