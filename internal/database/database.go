package database

import (
	"github.com/chaisql/chai"
)

func SetupDatabase() (*chai.DB, error) {
	db, err := chai.Open(":memory:")
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS brags (
			id TEXT PRIMARY KEY,
			description TEXT NOT NULL,
			details TEXT,
			created_at INTEGER NOT NULL,  -- Store as Unix timestamp
			updated_at INTEGER            -- Store as Unix timestamp
		)
	`)
	if err != nil {
		return nil, err
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
