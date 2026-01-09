package db

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// User Deposit
type NewDepositParams struct {
	DepositAmount string
	DepositDate   string
}

type UserDeposit struct {
	Value       int
	DepositDate time.Time
}
type Deposit struct {
	Id                       int       `db:"id"`
	DepositDate              time.Time `db:"deposit_date"`
	DepositAmountInEurocents int       `db:"deposit_amount_in_eurocents"`
	deposit_volume           float32
}

func (d *UserDeposit) From(depositParams NewDepositParams) error {
	if err := d.initValue(depositParams.DepositAmount); err != nil {
		return err
	}

	if err := d.initDate(depositParams.DepositDate); err != nil {
		return err
	}

	return nil
}
func (d *UserDeposit) initValue(depositAmountInput string) error {
	// @TODO - Handle numbers with precition greater than 2 decimals
	depositAmountInput = strings.Replace(depositAmountInput, ",", ".", 1)

	splittedUserInput := strings.Split(depositAmountInput, ".")

	if len(splittedUserInput) < 2 {
		splittedUserInput = []string{depositAmountInput, "0"}
	}

	integerPart, fractionalPart := splittedUserInput[0], splittedUserInput[1]
	fractionalPartLength := len(fractionalPart)

	parsedIntegerPart, err := strconv.Atoi(integerPart)
	parsedFractionalPart, err := strconv.Atoi(fractionalPart)

	if err != nil {
		return err
	}

	d.Value = parsedIntegerPart*int(math.Pow(10, float64(fractionalPartLength))) + parsedFractionalPart
	return nil
}
func (d *UserDeposit) initDate(depositDateInput string) error {
	if depositDateInput == "" {
		d.DepositDate = time.Time{}
		return nil
	}

	possibleFormats := []string{"02.01.2006", "02-01-2006", time.DateOnly}

	for _, format := range possibleFormats {

		parsedTime, err := time.Parse(format, depositDateInput)
		if err != nil {
			continue
		}

		d.DepositDate = parsedTime
		return nil
	}

	return nil
}

type UserRate struct {
	Date        time.Time
	RateInGrosz int
}

type NewRateParams struct {
	Date string
	Rate string
}

func (d *UserRate) From(params NewRateParams) error {
	parsedDate, err := time.Parse(time.DateOnly, params.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate

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

type UserIndexPrice struct {
	Date             time.Time
	PriceInEurocents int
}

type NewIndexPriceParams struct {
	Date             string
	PriceInEurocents string
}

func (d *UserIndexPrice) From(params NewIndexPriceParams) error {
	parsedDate, err := time.Parse(time.DateOnly, params.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate

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
