package actions

import (
	"errors"
	"fmt"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
)

func HandleUserAction() error {

	userInput, err := ParseUserInput()

	if err != nil {
		return err
	}

	switch userInput.ActionName {
	case "add":
		if err := handleSync(); err != nil {
			return err
		}

		if err := handleAddDeposit(userInput.ActionParams); err != nil {
			return err
		}
		return nil

	case "sync":
		if err := handleSync(); err != nil {
			return err
		}
		return nil

	case "status":
		if err := handleSync(); err != nil {
			return err
		}

		if err := handleStatus(); err != nil {
			return err
		}
		return nil

	case "help":
		ui.DisplayHelpScreen()

	default:
		return nil
	}

	return nil
}

func handleAddDeposit(addDepositParams []string) error {
	if len(addDepositParams) == 0 {
		return errors.New("Not enough parameters passed")
	}

	parsedDepositParams := db.NewDepositParams{
		DepositAmount: addDepositParams[0],
	}

	if len(addDepositParams) == 2 {
		parsedDepositParams.DepositDate = addDepositParams[1]
	}

	deposit := db.UserDeposit{}

	if err := deposit.From(parsedDepositParams); err != nil {
		return err
	}

	if err := db.AddDeposit(deposit); err != nil {
		return err
	}

	fmt.Println("Deposit added!")

	return nil
}

func handleSync() error {
	downloaders := []Downloader{
		NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api"),
		NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com"),
	}

	for _, downloader := range downloaders {
		err := downloader.SyncData()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func handleStatus() error {
	overallAmountInEurocents, err := db.GetOverallDepositInEurocents()
	if err != nil {
		return err
	}

	latestIndexPrice, err := db.GetLatestIndexPrice()
	if err != nil {
		return err
	}

	latestExchangeRate, err := db.GetLatestExchangeRate()
	if err != nil {
		fmt.Println("Warning: Could not fetch exchange rates.", err)
	}

	allDeposits, err := db.GetAllDeposits()
	if err != nil {
		return err
	}

	var totalUnitsOwned float64
	for _, deposit := range allDeposits {
		historicalPrice, err := db.GetIndexPriceByDate(deposit.DepositDate)
		if err != nil {
			fmt.Printf("Warning: Could not find index price for date %s. Skipping calculation for this deposit.\n", deposit.DepositDate)
			continue
		}

		unitsBought := float64(deposit.DepositAmountInEurocents) / float64(historicalPrice.PriceInEurocents)
		totalUnitsOwned += unitsBought
	}

	currentValue := totalUnitsOwned * float64(latestIndexPrice.PriceInEurocents)

	profitValue := currentValue - float64(overallAmountInEurocents)
	profitPercentage := (profitValue / float64(overallAmountInEurocents)) * 100

	fmt.Printf("Total Invested: %.2f EUR\n", float64(overallAmountInEurocents)/100)
	fmt.Printf("Current Value:  %.2f EUR\n", currentValue/100)
	fmt.Printf("Profit/Loss:    %.2f EUR (%.2f%%)\n", profitValue/100, profitPercentage)

	if latestExchangeRate.RateInGrosz > 0 {
		rateVal := float64(latestExchangeRate.RateInGrosz) / 100.0
		currentValuePLN := (currentValue / 100) * rateVal
		profitValuePLN := (profitValue / 100) * rateVal

		fmt.Println("---------------------------")
		fmt.Printf("Rate (1 EUR):   %.4f PLN\n", rateVal)
		fmt.Printf("Assets Value:   %.2f PLN\n", currentValuePLN)
		fmt.Printf("Profit/Loss: %.2f PLN\n", profitValuePLN)
	}

	return nil
}
