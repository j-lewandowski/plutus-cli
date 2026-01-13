package db

import "time"

func (r *Repository) GetFirstDeposit() (Deposit, error) {
	data := r.conn.QueryRow(`
		SELECT * FROM deposit
		ORDER BY deposit_date
		LIMIT 1;
	`)

	firstDeposit := Deposit{}
	if err := data.Scan(
		&firstDeposit.Id,
		&firstDeposit.DepositDate,
		&firstDeposit.DepositAmountInEurocents,
		&firstDeposit.DepositVolume,
		&firstDeposit.DepositVolumePrecision); err != nil {
		return Deposit{}, err
	}

	return firstDeposit, nil
}

func (r *Repository) GetLastDeposit() (Deposit, error) {
	data := r.conn.QueryRow(`
		SELECT * FROM deposit
		ORDER BY deposit_date DESC
		LIMIT 1;
	`)

	lastDeposit := Deposit{}
	if err := data.Scan(
		&lastDeposit.Id,
		&lastDeposit.DepositDate,
		&lastDeposit.DepositAmountInEurocents,
		&lastDeposit.DepositVolume,
		&lastDeposit.DepositVolumePrecision); err != nil {
		return Deposit{}, err
	}

	return lastDeposit, nil
}

func (r *Repository) GetAllDeposits() ([]Deposit, error) {
	rows, err := r.conn.Query(`
        SELECT * FROM deposit
        ORDER BY deposit_date DESC;
    `)
	if err != nil {
		return []Deposit{}, err
	}
	defer rows.Close()

	var deposits []Deposit

	for rows.Next() {
		var d Deposit
		if err := rows.Scan(
			&d.Id,
			&d.DepositDate,
			&d.DepositAmountInEurocents,
			&d.DepositVolume,
			&d.DepositVolumePrecision,
		); err != nil {
			return []Deposit{}, err
		}
		deposits = append(deposits, d)
	}

	return deposits, nil
}

func (r *Repository) AddDeposit(deposit UserDeposit) error {
	_, err := r.conn.Exec(`
	INSERT INTO deposit (deposit_date, deposit_amount_in_eurocents, deposit_volume, deposit_volume_precision)
	VALUES ($1, $2, $3, $4);`, deposit.DepositDate, deposit.Value, deposit.Volume, deposit.VolumePrecision)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddRates(rates []CurrencyRate) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO eur_exchange_rate (date, price_pln_in_grosz) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rate := range rates {
		if _, err := stmt.Exec(rate.Date, rate.RateInGrosz); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) AddIndexPrices(indexPrices []IndexPrice) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO index_price (date, price_in_eurocents) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rate := range indexPrices {
		if _, err := stmt.Exec(rate.Date, rate.PriceInEurocents); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetOverallDepositInEurocents() (int, error) {
	data := r.conn.QueryRow(`
		SELECT SUM(deposit_amount_in_eurocents) FROM deposit;
	`)

	var overallDeposit int = 0

	if err := data.Scan(&overallDeposit); err != nil {
		return 0, err
	}

	return overallDeposit, nil
}

func (r *Repository) GetLatestIndexPrice() (IndexPrice, error) {
	data := r.conn.QueryRow(`
    SELECT * FROM index_price
    ORDER BY date DESC
    LIMIT 1;`)

	latestIndexPrice := IndexPrice{}
	if err := data.Scan(
		&latestIndexPrice.Date,
		&latestIndexPrice.PriceInEurocents); err != nil {
		return IndexPrice{}, err
	}

	return latestIndexPrice, nil
}

func (r *Repository) GetIndexPriceByDate(date time.Time) (IndexPrice, error) {
	row := r.conn.QueryRow(`
            SELECT date, price_in_eurocents 
            FROM index_price 
            WHERE date <= ?
            ORDER BY date DESC
            LIMIT 1;
    `, date)

	var indexPrice IndexPrice
	if err := row.Scan(
		&indexPrice.Date,
		&indexPrice.PriceInEurocents,
	); err != nil {
		return IndexPrice{}, err
	}
	return indexPrice, nil
}

func (r *Repository) GetLatestExchangeRate() (CurrencyRate, error) {
	row := r.conn.QueryRow(`
        SELECT date, price_pln_in_grosz
        FROM eur_exchange_rate
        ORDER BY date DESC
        LIMIT 1;
    `)

	var rate CurrencyRate
	if err := row.Scan(&rate.Date, &rate.RateInGrosz); err != nil {
		return CurrencyRate{}, err
	}

	return rate, nil
}
