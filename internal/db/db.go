package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDb() (*sql.DB, error) {

	home, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(home, ".plutus.sqlite")


	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}