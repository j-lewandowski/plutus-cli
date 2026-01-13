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

func (d *UserDeposit) From(depositParams NewDepositParams) error {
	if err := d.initValue(depositParams.DepositAmount); err != nil {
		return err
	}

	if err := d.initVolume(depositParams.DepositVolume); err != nil {
		return err
	}

	if err := d.initDate(depositParams.DepositDate); err != nil {
		return err
	}

	return nil
}

func (d *UserDeposit) initVolume(depositVolumeInput string) error {
	if depositVolumeInput == "" {
		d.Volume = 0
		d.VolumePrecision = 0
		return nil
	}

	depositVolumeInput = strings.Replace(depositVolumeInput, ",", ".", 1)
	splittedUserInput := strings.Split(depositVolumeInput, ".")

	if len(splittedUserInput) < 2 {
		val, err := strconv.ParseInt(depositVolumeInput, 10, 64)
		if err != nil {
			return err
		}
		d.Volume = val
		d.VolumePrecision = 0
		return nil
	}

	integerPart, fractionalPart := splittedUserInput[0], splittedUserInput[1]
	fractionalPartLength := len(fractionalPart)

	// Combine integer and fractional parts
	combined := integerPart + fractionalPart
	val, err := strconv.ParseInt(combined, 10, 64)
	if err != nil {
		return err
	}

	d.Volume = val
	d.VolumePrecision = fractionalPartLength
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

type IndexPrice struct {
	Date             time.Time
	PriceInEurocents int
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
