package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{conn: db}
}

func InitDb() (*Repository, error) {

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

	repo := NewRepository(db)
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) Close() {
	r.conn.Close()
}

func (r *Repository) migrate() error {
	_, err := r.conn.Exec(`
    CREATE TABLE IF NOT EXISTS "deposit" (
        id 													INTEGER 					PRIMARY KEY 																AUTOINCREMENT,
        deposit_date								DATE							DEFAULT(datetime(current_timestamp)),
        deposit_amount_in_eurocents	INTEGER						NOT NULL,
        deposit_volume							INTEGER						NOT NULL,
        deposit_volume_precision		INTEGER						NOT NULL
    );`)

	if err != nil {
		return err
	}

	_, err = r.conn.Exec(`
    CREATE TABLE IF NOT EXISTS "index_price" (
        date										DATE						PRIMARY KEY 	NOT NULL,
        price_in_eurocents			INTEGER					NOT NULL
    );`)

	if err != nil {
		return err
	}

	_, err = r.conn.Exec(`
    CREATE TABLE IF NOT EXISTS "eur_exchange_rate" (
        date									DATE						PRIMARY KEY 	NOT NULL,
        price_pln_in_grosz		INTEGER					NOT NULL
    );`)

	if err != nil {
		return err
	}
	return nil
}
