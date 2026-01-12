package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var connection *sql.DB

func InitDb() error {

	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	dbPath := filepath.Join(home, ".plutus.sqlite")

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "deposit" (
		id 													INTEGER 					PRIMARY KEY 																AUTOINCREMENT,
		deposit_date								DATE							DEFAULT(datetime(current_timestamp)),
		deposit_amount_in_eurocents	INTEGER						NOT NULL,
		deposit_volume							INTEGER				NOT NULL
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "index_price" (
		date										DATE						PRIMARY KEY 	NOT NULL,
		price_in_eurocents			INTEGER					NOT NULL
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "eur_exchange_rate" (
		date									DATE						PRIMARY KEY 	NOT NULL,
		price_pln_in_grosz		INTEGER					NOT NULL
	);`)

	if err != nil {
		fmt.Println("Here")
		return err
	}

	connection = db
	return nil
}

func GetDb() *sql.DB {
	return connection
}

func Close() {
	connection.Close()
}
