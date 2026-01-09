package db

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
		&firstDeposit.deposit_volume); err != nil {
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
		&lastDeposit.deposit_volume); err != nil {
		return Deposit{}, err
	}

	return lastDeposit, nil
}

func AddDeposit(deposit UserDeposit) error {
	db := GetDb()

	_, err := db.Exec(`
	INSERT INTO deposit (deposit_date, deposit_amount_in_eurocents, deposit_volume)
	VALUES ($1, $2, $3);`, deposit.DepositDate, deposit.Value, 0)

	if err != nil {
		return err
	}

	return nil
}

func AddRates(rates []UserRate) error {
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

func AddIndexPrices(indexPrices []UserIndexPrice) error {
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
