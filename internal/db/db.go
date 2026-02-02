package db

import (
	"database/sql"
	"fmt"
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
	var dbPath string

	envPath := os.Getenv("PLUTUS_DB")

	if envPath != "" {
		dbPath = envPath
		fmt.Printf("Using database path from PLUTUS_DB: %s\n", dbPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(home, ".plutus.sqlite")
	}

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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS "deposit" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			deposit_date DATE DEFAULT(datetime(current_timestamp)),
			deposit_amount_in_eurocents INTEGER NOT NULL,
			deposit_volume INTEGER NOT NULL,
			deposit_volume_precision INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS "eur_exchange_rate" (
			date DATE PRIMARY KEY NOT NULL,
			price_pln_in_grosz INTEGER NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := r.conn.Exec(q); err != nil {
			return err
		}
	}

	var columnExists bool
	err := r.conn.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pragma_table_info('index_price') WHERE name='is_real'
		);
	`).Scan(&columnExists)

	if err != nil || !columnExists {
		fmt.Println("🔄 Database migration: Rebuilding index_price table...")

		r.conn.Exec(`DROP TABLE IF EXISTS index_price;`)
		_, err = r.conn.Exec(`
			CREATE TABLE index_price (
				date               DATE    PRIMARY KEY NOT NULL,
				price_in_eurocents INTEGER NOT NULL,
				is_real            BOOLEAN NOT NULL DEFAULT 0
			);`)
		if err != nil {
			return fmt.Errorf("failed to recreate index_price: %w", err)
		}
	}

	return nil
}
