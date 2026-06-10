package db

import (
	"database/sql"
	"fmt"
	"time"
)

func (r *Repository) GetFirstDeposit() (Deposit, error) {
	data := r.conn.QueryRow(`
		SELECT * FROM deposit
		ORDER BY deposit_date
		LIMIT 1;
	`)

	firstDeposit := Deposit{}
	err := data.Scan(
		&firstDeposit.Id,
		&firstDeposit.DepositDate,
		&firstDeposit.DepositAmountInEurocents,
		&firstDeposit.DepositVolume,
		&firstDeposit.DepositVolumePrecision)

	if err == sql.ErrNoRows {
		return Deposit{}, nil
	}
	if err != nil {
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
	VALUES ($1, $2, $3, $4);`, deposit.DepositDate.Format("2006-01-02"), deposit.Value, deposit.Volume, deposit.VolumePrecision)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllSells() ([]Sell, error) {
	rows, err := r.conn.Query(`
        SELECT * FROM sell
        ORDER BY sell_date DESC;
  `)
	if err != nil {
		return []Sell{}, err
	}
	defer rows.Close()

	var sells []Sell

	for rows.Next() {
		var s Sell
		if err := rows.Scan(
			&s.Id,
			&s.SellDate,
			&s.SellAmountInEurocents,
			&s.SellVolume,
			&s.SellVolumePrecision,
		); err != nil {
			return []Sell{}, err
		}
		sells = append(sells, s)
	}

	return sells, nil
}

func (r *Repository) AddSell(sell UserSell) error {
	_, err := r.conn.Exec(`
	INSERT INTO sell (sell_date, sell_amount_in_eurocents, sell_volume, sell_volume_precision)
	VALUES ($1, $2, $3, $4);`, sell.SellDate.Format("2006-01-02"), sell.Value, sell.Volume, sell.VolumePrecision)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetOverallSellInEurocents() (int, error) {
	data := r.conn.QueryRow(`
		SELECT COALESCE(SUM(sell_amount_in_eurocents), 0) FROM sell;
	`)

	var overallSell int = 0

	if err := data.Scan(&overallSell); err != nil {
		return 0, err
	}

	return overallSell, nil
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
		if _, err := stmt.Exec(rate.Date.Format("2006-01-02"), rate.RateInGrosz); err != nil {
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

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO index_price (date, price_in_eurocents, is_real) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range indexPrices {
		if _, err := stmt.Exec(p.Date.Format("2006-01-02"), p.PriceInEurocents, p.IsReal); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) GetOverallDepositInEurocents() (int, error) {
	data := r.conn.QueryRow(`
		SELECT COALESCE(SUM(deposit_amount_in_eurocents), 0) FROM deposit;
	`)

	var overallDeposit int = 0

	if err := data.Scan(&overallDeposit); err != nil {
		return 0, err
	}

	return overallDeposit, nil
}

func (r *Repository) GetLatestIndexPrice() (IndexPrice, error) {
	data := r.conn.QueryRow(`
    SELECT date, price_in_eurocents, is_real 
    FROM index_price
    ORDER BY date DESC
    LIMIT 1;`)

	latestIndexPrice := IndexPrice{}
	if err := data.Scan(
		&latestIndexPrice.Date,
		&latestIndexPrice.PriceInEurocents,
		&latestIndexPrice.IsReal,
	); err != nil {
		return IndexPrice{}, err
	}

	return latestIndexPrice, nil
}

func (r *Repository) GetIndexPriceByDate(date time.Time) (IndexPrice, error) {
	row := r.conn.QueryRow(`
			SELECT date, price_in_eurocents, is_real 
			FROM index_price 
			WHERE date >= ? AND is_real = 1
			ORDER BY date ASC 
			LIMIT 1;
	`, date.Format("2006-01-02"))

	var indexPrice IndexPrice
	if err := row.Scan(
		&indexPrice.Date,
		&indexPrice.PriceInEurocents,
		&indexPrice.IsReal,
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

func (r *Repository) GetChartPoints() ([]ChartPoint, error) {
	query := `
	WITH daily_purchases AS (
		SELECT 
			d.deposit_date as d_date,
			SUM(d.deposit_amount_in_eurocents) AS daily_invested,
			SUM(CAST(d.deposit_amount_in_eurocents AS FLOAT) / ip.price_in_eurocents) AS units_bought
		FROM deposit d
		JOIN index_price ip ON d.deposit_date = ip.date
		GROUP BY d_date
	),
	daily_sells AS (
		SELECT 
			s.sell_date as d_date,
			SUM(s.sell_amount_in_eurocents) AS daily_withdrawn,
			SUM(CAST(s.sell_amount_in_eurocents AS FLOAT) / ip.price_in_eurocents) AS units_sold
		FROM sell s
		JOIN index_price ip ON s.sell_date = ip.date
		GROUP BY d_date
	),
	running_totals AS (
		SELECT 
			ip.date as d_date,
			SUM(COALESCE(dp.daily_invested, 0)) OVER (ORDER BY ip.date)
				- SUM(COALESCE(ds.daily_withdrawn, 0)) OVER (ORDER BY ip.date) AS total_invested,
			SUM(COALESCE(dp.units_bought, 0)) OVER (ORDER BY ip.date)
				- SUM(COALESCE(ds.units_sold, 0)) OVER (ORDER BY ip.date) AS total_units,
			ip.price_in_eurocents
		FROM index_price ip
		LEFT JOIN daily_purchases dp ON ip.date = dp.d_date
		LEFT JOIN daily_sells ds ON ip.date = ds.d_date
		WHERE ip.date >= (
			SELECT MIN(event_date) FROM (
				SELECT deposit_date AS event_date FROM deposit
				UNION ALL
				SELECT sell_date AS event_date FROM sell
			)
		)
	)
	SELECT 
		d_date,
		CAST(total_invested AS FLOAT) / 100.0 AS Invested,
		(total_units * price_in_eurocents) / 100.0 AS Value
	FROM running_totals
	WHERE total_invested > 0 OR total_units > 0
	ORDER BY d_date;`

	rows, err := r.conn.Query(query)
	if err != nil {
		return []ChartPoint{}, fmt.Errorf("SQL Error: %w", err)
	}
	defer rows.Close()

	var chartPoints []ChartPoint
	for rows.Next() {
		var d ChartPoint
		if err := rows.Scan(&d.Date, &d.Invested, &d.Value); err != nil {
			return []ChartPoint{}, err
		}
		chartPoints = append(chartPoints, d)
	}

	return chartPoints, nil
}
