package actions

import (
	"fmt"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
	"plutus-cli/internal/portfolio"
	"sync"
)

type Handler struct {
	Repo *db.Repository
}

func NewHandler(repo *db.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

func (h *Handler) Run() error {

	userInput, err := ParseUserInput()

	if err != nil {
		return err
	}

	switch userInput.ActionName {
	case "add":
		if err := h.handleSync(); err != nil {
			return err
		}

		if err := h.handleAddDeposit(userInput.ActionParams); err != nil {
			return err
		}
		return nil

	case "sync":
		if err := h.handleSync(); err != nil {
			return err
		}
		return nil

	case "status":
		if err := h.handleSync(); err != nil {
			return err
		}

		if err := h.handleStatus(); err != nil {
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

func (h *Handler) handleAddDeposit(addDepositParams []string) error {
	if len(addDepositParams) == 0 {
		return fmt.Errorf("Not enough parameters passed")
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

	if err := h.Repo.AddDeposit(deposit); err != nil {
		return err
	}

	fmt.Println("Deposit added!")

	return nil
}

func (h *Handler) handleSync() error {
	downloaders := []Downloader{
		NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api", h.Repo),
		NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com", h.Repo),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(downloaders))

	for _, downloader := range downloaders {
		wg.Add(1)

		go func(d Downloader) {
			defer wg.Done()
			if err := d.SyncData(); err != nil {
				errChan <- fmt.Errorf("%s failed: %w", d.GetName(), err)
			}
		}(downloader)
	}

	wg.Wait()
	close(errChan)

	var combinedError error
	for err := range errChan {
		if combinedError == nil {
			combinedError = err
		} else {
			fmt.Println("Error:", err)
		}
	}

	if combinedError != nil {
		return fmt.Errorf("sync completed with errors")
	}

	fmt.Println("Sync completed successfully.")
	return nil
}

func (h *Handler) handleStatus() error {
	report, err := portfolio.CalculatePortfolio(h.Repo)
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
