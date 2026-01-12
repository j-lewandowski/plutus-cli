package db

import "time"

func GetFirstDeposit() (Deposit, error) {
	db := GetDb()

	data := db.QueryRow(`
		SELECT * FROM deposit
		ORDER BY deposit_date
		LIMIT 1;
	`)

	firstDeposit := Deposit{}
	if err := data.Scan(
		&firstDeposit.Id,
		&firstDeposit.DepositDate,
		&firstDeposit.DepositAmountInEurocents,
		&firstDeposit.DepositVolume); err != nil {
		return Deposit{}, err
	}

	return firstDeposit, nil
}

func GetLastDeposit() (Deposit, error) {
	db := GetDb()

	data := db.QueryRow(`
		SELECT * FROM deposit
		ORDER BY deposit_date DESC
		LIMIT 1;
	`)

	lastDeposit := Deposit{}
	if err := data.Scan(
		&lastDeposit.Id,
		&lastDeposit.DepositDate,
		&lastDeposit.DepositAmountInEurocents,
		&lastDeposit.DepositVolume); err != nil {
		return Deposit{}, err
	}

	return lastDeposit, nil
}

func GetAllDeposits() ([]Deposit, error) {
	db := GetDb()

	rows, err := db.Query(`
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
		); err != nil {
			return []Deposit{}, err
		}
		deposits = append(deposits, d)
	}

	return deposits, nil
}

func AddDeposit(deposit UserDeposit) error {
	db := GetDb()

	_, err := db.Exec(`
	INSERT INTO deposit (deposit_date, deposit_amount_in_eurocents, deposit_volume)
	VALUES ($1, $2, $3);`, deposit.DepositDate, deposit.Value, deposit.Volume)

	if err != nil {
		return err
	}

	return nil
}

func AddRates(rates []CurrencyRate) error {
	db := GetDb()

	sqlStr := "INSERT OR IGNORE INTO eur_exchange_rate (date, price_pln_in_grosz) VALUES "
	values := []interface{}{}

	for _, rate := range rates {
		sqlStr += "(?, ?),"
		values = append(values, rate.Date, rate.RateInGrosz)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)

	if err != nil {
		return err
	}

	return nil
}

func AddIndexPrices(indexPrices []IndexPrice) error {
	db := GetDb()

	sqlStr := "INSERT OR IGNORE INTO index_price (date, price_in_eurocents) VALUES "
	values := []interface{}{}

	for _, rate := range indexPrices {
		sqlStr += "(?, ?),"
		values = append(values, rate.Date, rate.PriceInEurocents)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)

	if err != nil {
		return err
	}

	return nil
}

func GetOverallDepositInEurocents() (int, error) {
	db := GetDb()

	data := db.QueryRow(`
		SELECT SUM(deposit_amount_in_eurocents) FROM deposit;
	`)

	var overallDeposit int = 0

	if err := data.Scan(&overallDeposit); err != nil {
		return 0, err
	}

	return overallDeposit, nil
}

func GetLatestIndexPrice() (IndexPrice, error) {
	db := GetDb()

	data := db.QueryRow(`
    SELECT * FROM index_price
    ORDER BY date DESC
    LIMIT 1;`)

	lastestIndexPrice := IndexPrice{}
	if err := data.Scan(
		&lastestIndexPrice.Date,
		&lastestIndexPrice.PriceInEurocents); err != nil {
		return IndexPrice{}, err
	}

	return lastestIndexPrice, nil
}

func GetIndexPriceByDate(date time.Time) (IndexPrice, error) {
	db := GetDb()

	row := db.QueryRow(`
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

func GetLatestExchangeRate() (CurrencyRate, error) {
	db := GetDb()

	row := db.QueryRow(`
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
