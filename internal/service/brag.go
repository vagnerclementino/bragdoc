package service

import (
	"database/sql"
	"fmt"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/usercase"
	"log"

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

	// Insert the Brag instance into the database
	insertQuery := `
		INSERT INTO brags (id, description, details, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?);
	`

	_, err := s.db.Exec(insertQuery, brag.ID, brag.Description, brag.Details, brag.CreatedAt, brag.UpdatedAt)
	if err != nil {
		return err
	}

	fmt.Println("Brag instance saved to the database.")
	return nil
}
func NewBragService() usercase.BragUserCase {
	db, err := sql.Open("sqlite3", "brags.db")
	// Create the Brags table if it doesn't exist
	createTable := `
		CREATE TABLE IF NOT EXISTS brags (
			id TEXT PRIMARY KEY,
			description TEXT,
			details TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);
	`
	_, err = db.Exec(createTable)

	if err != nil {
		log.Fatal(err)
	}
	return &bragService{
		db: db,
	}
}
