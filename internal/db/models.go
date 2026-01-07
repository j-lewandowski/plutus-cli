package db

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// User Deposit
type AddDepositParams struct {
	DepositAmount string
	DepositDate   string
}

type UserDeposit struct {
	Value       int
	DepositDate time.Time
}

func (d *UserDeposit) From(depositParams AddDepositParams) error {
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
