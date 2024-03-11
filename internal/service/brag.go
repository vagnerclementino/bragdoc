package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/usercase"

	_ "github.com/mattn/go-sqlite3"
)

type bragService struct {
	db *sql.DB
}

func (s *bragService) AddBrag(brag *domain.Brag) error {

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s.db)

	insertQuery := `
		INSERT INTO brags (id, description, details, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?);`

	_, err := s.db.Exec(insertQuery, brag.ID, brag.Description, brag.Details, brag.CreatedAt, brag.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
func NewBragService() usercase.BragUserCase {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	hiddenFolderPath := filepath.Join(homeDir, ".bragdoc")

	err = os.MkdirAll(hiddenFolderPath, 0700)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/brags.db", hiddenFolderPath))
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS brags (
			id TEXT PRIMARY KEY,
			description TEXT,
			details TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);`
	_, err = db.Exec(createTable)

	if err != nil {
		log.Fatal(err)
	}
	return &bragService{
		db: db,
	}
}
