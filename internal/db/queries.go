package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

func AddDeposit(deposit UserDeposit) error {
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	dbPath := filepath.Join(home, ".plutus.sqlite")

	db, err := sql.Open("sqlite", dbPath)

	_, err = db.Exec(`
	INSERT INTO deposit (deposit_date, deposit_amount_in_eurocents, deposit_volume)
	VALUES ($1, $2, $3);`, deposit.DepositDate, deposit.Value, 0)

	if err != nil {
		return err
	}

	return nil
}
