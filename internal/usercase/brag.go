package usercase

import "github.com/vagnerclementino/bragdoc/internal/domain"

type BragUserCase interface {
	AddBrag(brag *domain.Brag) error
}
