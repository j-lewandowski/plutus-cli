package actions

import (
	"errors"
	"fmt"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
	"plutus-cli/internal/portfolio"
)

func HandleUserAction(repo *db.Repository) error {

	userInput, err := ParseUserInput()

	if err != nil {
		return err
	}

	switch userInput.ActionName {
	case "add":
		if err := handleSync(repo); err != nil {
			return err
		}

		if err := handleAddDeposit(repo, userInput.ActionParams); err != nil {
			return err
		}
		return nil

	case "sync":
		if err := handleSync(repo); err != nil {
			return err
		}
		return nil

	case "status":
		if err := handleSync(repo); err != nil {
			return err
		}

		if err := handleStatus(repo); err != nil {
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

func handleAddDeposit(repo *db.Repository, addDepositParams []string) error {
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

	if err := repo.AddDeposit(deposit); err != nil {
		return err
	}

	fmt.Println("Deposit added!")

	return nil
}

func handleSync(repo *db.Repository) error {
	downloaders := []Downloader{
		NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api", repo),
		NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com", repo),
	}

	for _, downloader := range downloaders {
		err := downloader.SyncData()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func handleStatus(repo *db.Repository) error {
	report, err := portfolio.CalculatePortfolio(repo)
	if err != nil {
		return err
	}

	for _, warning := range report.Warnings {
		fmt.Println(warning)
	}

	fmt.Printf("Total Invested: %.2f EUR\n", float64(report.TotalInvestedInEurocents)/100.0)
	fmt.Printf("Current Value:  %.2f EUR\n", float64(report.CurrentValueInEurocents)/100.0)
	fmt.Printf("Profit/Loss:    %.2f EUR (%.2f%%)\n", float64(report.ProfitValueInEurocents)/100.0, report.ProfitPercent)

	if report.HasExchangeRate {
		fmt.Println("---------------------------")
		fmt.Printf("Rate (1 EUR):   %.4f PLN\n", float64(report.RateEURtoPLNInGrosz)/100.0)
		fmt.Printf("Assets Value:   %.2f PLN\n", float64(report.CurrentValueInGrosz)/100.0)
		fmt.Printf("Profit/Loss: %.2f PLN\n", float64(report.ProfitValueInGrosz)/100.0)
	}

	return nil
}
