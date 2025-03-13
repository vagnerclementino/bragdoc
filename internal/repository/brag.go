package repository

import (
	"context"

	"github.com/chaisql/chai"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

type BragRepository interface {
	Create(ctx context.Context, brag *domain.Brag) error
	FindAll(ctx context.Context) ([]domain.Brag, error)
}

type bragRepo struct {
	db *chai.DB
}

func NewBragRepository(db *chai.DB) BragRepository {
	return &bragRepo{db: db}
}

func (r *bragRepo) Create(ctx context.Context, brag *domain.Brag) error {
	return r.db.Exec(`
		INSERT INTO brags (id, description, details, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`,
		brag.ID, brag.Description, brag.Details, brag.CreatedAt, brag.UpdatedAt)
}

func (r *bragRepo) FindAll(ctx context.Context) ([]domain.Brag, error) {
	var brags []domain.Brag

	res, err := r.db.Query("SELECT * FROM brags")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	err = res.Iterate(func(r *chai.Row) error {
		var b domain.Brag
		var createdAt, updatedAt int64
		if err := r.Scan(&b.ID, &b.Description, &b.Details, &createdAt, &updatedAt); err != nil {
			return err
		}
		b.CreatedAt = createdAt
		if updatedAt != 0 {
			b.UpdatedAt = &updatedAt
		}
		brags = append(brags, b)
		return nil
	})

	return brags, err
}
