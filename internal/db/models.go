package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NewDepositParams struct {
	DepositAmount string
	DepositVolume string
	DepositDate   string
}

type UserDeposit struct {
	Value           int
	Volume          int64
	VolumePrecision int
	DepositDate     time.Time
}
type Deposit struct {
	Id                       int       `db:"id"`
	DepositDate              time.Time `db:"deposit_date"`
	DepositAmountInEurocents int       `db:"deposit_amount_in_eurocents"`
	DepositVolume            int64     `db:"deposit_volume"`
	DepositVolumePrecision   int       `db:"deposit_volume_precision"`
}

type NewSellParams struct {
	SellAmount string
	SellVolume string
	SellDate   string
}

type UserSell struct {
	Value           int
	Volume          int64
	VolumePrecision int
	SellDate        time.Time
}

type Sell struct {
	Id                    int       `db:"id"`
	SellDate              time.Time `db:"sell_date"`
	SellAmountInEurocents int       `db:"sell_amount_in_eurocents"`
	SellVolume            int64     `db:"sell_volume"`
	SellVolumePrecision   int       `db:"sell_volume_precision"`
}

func (s *UserSell) From(sellParams NewSellParams) error {
	value, err := parseEurocents(sellParams.SellAmount)
	if err != nil {
		return err
	}
	s.Value = value

	volume, precision, err := parseOptionalVolume(sellParams.SellVolume)
	if err != nil {
		return err
	}
	s.Volume = volume
	s.VolumePrecision = precision

	sellDate, err := parseUserDate(sellParams.SellDate)
	if err != nil {
		return err
	}
	s.SellDate = sellDate

	return nil
}

func (d *UserDeposit) From(depositParams NewDepositParams) error {
	value, err := parseEurocents(depositParams.DepositAmount)
	if err != nil {
		return err
	}
	d.Value = value

	volume, precision, err := parseOptionalVolume(depositParams.DepositVolume)
	if err != nil {
		return err
	}
	d.Volume = volume
	d.VolumePrecision = precision

	depositDate, err := parseUserDate(depositParams.DepositDate)
	if err != nil {
		return err
	}
	d.DepositDate = depositDate

	return nil
}

func parseOptionalVolume(volumeInput string) (int64, int, error) {
	if volumeInput == "" {
		return 0, 0, nil
	}

	volumeInput = strings.Replace(volumeInput, ",", ".", 1)
	splittedUserInput := strings.Split(volumeInput, ".")

	if len(splittedUserInput) < 2 {
		val, err := strconv.ParseInt(volumeInput, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return val, 0, nil
	}

	integerPart, fractionalPart := splittedUserInput[0], splittedUserInput[1]
	combined := integerPart + fractionalPart
	val, err := strconv.ParseInt(combined, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return val, len(fractionalPart), nil
}

func parseEurocents(amountInput string) (int, error) {
	amountInput = strings.Replace(amountInput, ",", ".", 1)

	parts := strings.Split(amountInput, ".")

	if len(parts) > 1 && len(parts[1]) > 2 {
		return 0, fmt.Errorf("Deposit amount cannot have more than 2 decimal places")
	}

	integerPart := parts[0]
	if integerPart == "" {
		integerPart = "0"
	}

	fractionalPart := "00"
	if len(parts) > 1 {
		fractionalPart = parts[1]
		if len(fractionalPart) == 1 {
			fractionalPart = fractionalPart + "0"
		}
	}

	combined := integerPart + fractionalPart
	val, err := strconv.Atoi(combined)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func parseUserDate(dateInput string) (time.Time, error) {
	if dateInput == "" {
		return time.Now().UTC(), nil
	}

	normalized := strings.ReplaceAll(dateInput, ".", "-")

	const layout = "02-01-2006"
	parsedTime, err := time.Parse(layout, normalized)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid date format: '%s'. Use DD.MM.YYYY or DD-MM-YYYY", dateInput)
	}

	return parsedTime.UTC(), nil
}

type CurrencyRate struct {
	Date        time.Time
	RateInGrosz int
}

type NewRateParams struct {
	Date string
	Rate string
}

func (d *CurrencyRate) From(params NewRateParams) error {
	parsedDate, err := time.Parse(time.DateOnly, params.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate.UTC()

	rateStr := strings.Replace(params.Rate, ",", ".", 1)
	parts := strings.Split(rateStr, ".")

	whole := parts[0]
	fraction := "00"
	if len(parts) > 1 {
		fraction = parts[1]
		if len(fraction) >= 2 {
			fraction = fraction[:2]
		} else {
			fraction = fraction + "0"
		}
	}

	combined := whole + fraction
	val, err := strconv.Atoi(combined)
	if err != nil {
		return err
	}

	d.RateInGrosz = val
	return nil
}

type IndexPrice struct {
	Date             time.Time
	PriceInEurocents int
	IsReal           bool
}

type NewIndexPriceParams struct {
	Date             string
	PriceInEurocents string
}

func (d *IndexPrice) From(params NewIndexPriceParams) error {
	parsedDate, err := time.Parse(time.DateOnly, params.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate.UTC()

	rateStr := strings.Replace(params.PriceInEurocents, ",", ".", 1)
	parts := strings.Split(rateStr, ".")

	whole := parts[0]
	fraction := "00"
	if len(parts) > 1 {
		fraction = parts[1]
		if len(fraction) >= 2 {
			fraction = fraction[:2]
		} else {
			fraction = fraction + "0"
		}
	}

	combined := whole + fraction
	val, err := strconv.Atoi(combined)
	if err != nil {
		return err
	}

	d.PriceInEurocents = val
	return nil
}

func (c CurrencyRate) GetDate() time.Time {
	return c.Date
}

func (c CurrencyRate) CreateWithDate(date time.Time) interface{} {
	return CurrencyRate{
		Date:        date,
		RateInGrosz: c.RateInGrosz,
	}
}

func (i IndexPrice) GetDate() time.Time {
	return i.Date
}

func (i IndexPrice) CreateWithDate(date time.Time) interface{} {
	return IndexPrice{
		Date:             date,
		PriceInEurocents: i.PriceInEurocents,
	}
}

type ChartPoint struct {
	Date     string
	Invested float64
	Value    float64
}
