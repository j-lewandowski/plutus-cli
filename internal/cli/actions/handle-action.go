package actions

import (
	"fmt"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
	"plutus-cli/internal/portfolio"
	"sync"
)

type Handler struct {
	Repo             *db.Repository
	AvailableActions []Action
}

type Action interface {
	Run(args []string) error
	GetInputTrigger() string
	GetDescription() string
}

type ActionData struct {
	InputTrigger  string
	CanonicalName string
	Description   string
}

func (ad ActionData) GetInputTrigger() string {
	return ad.InputTrigger
}

func (ad ActionData) GetDescription() string {
	return ad.Description
}

type AddAction struct {
	ActionData
	Repo *db.Repository
}

func (a *AddAction) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not enough parameters passed")
	}

	parsedDepositParams := db.NewDepositParams{
		DepositAmount: args[0],
	}

	if len(args) == 2 {
		parsedDepositParams.DepositDate = args[1]
	}

	deposit := db.UserDeposit{}

	if err := deposit.From(parsedDepositParams); err != nil {
		return err
	}

	if err := a.Repo.AddDeposit(deposit); err != nil {
		return err
	}

	fmt.Println("Deposit added!")

	if err := runSync(a.Repo); err != nil {
		return err
	}

	return nil
}

type SyncAction struct {
	ActionData
	Repo *db.Repository
}

func (a *SyncAction) Run(args []string) error {
	return runSync(a.Repo)
}

type StatusAction struct {
	ActionData
	Repo *db.Repository
}

func (a *StatusAction) Run(args []string) error {
	if err := runSync(a.Repo); err != nil {
		return err
	}

	report, err := portfolio.CalculatePortfolio(a.Repo)
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

type HelpAction struct {
	ActionData
	AllActions []Action
}

func (a *HelpAction) Run(args []string) error {
	var entries []ui.HelpEntry
	for _, action := range a.AllActions {
		entries = append(entries, ui.HelpEntry{
			Name:        action.GetInputTrigger(),
			Description: action.GetDescription(),
		})
	}
	ui.DisplayHelpScreen(entries)
	return nil
}

func NewHandler(repo *db.Repository) *Handler {
	var actionsList []Action

	helpAction := &HelpAction{
		ActionData: ActionData{
			InputTrigger:  "help",
			CanonicalName: "help",
			Description:   "See more information on a command",
		},
	}

	actionsList = []Action{
		&AddAction{
			ActionData: ActionData{
				InputTrigger:  "add",
				CanonicalName: "add",
				Description:   "Allows user to add deposit event.",
			},
			Repo: repo,
		},
		&SyncAction{
			ActionData: ActionData{
				InputTrigger:  "sync",
				CanonicalName: "sync",
				Description:   "Syncs the CLI with up-to-date market data.",
			},
			Repo: repo,
		},
		&StatusAction{
			ActionData: ActionData{
				InputTrigger:  "status",
				CanonicalName: "status",
				Description:   "Displays current portfolio value and profit/loss percentage.",
			},
			Repo: repo,
		},
		helpAction,
	}

	helpAction.AllActions = actionsList

	return &Handler{
		Repo:             repo,
		AvailableActions: actionsList,
	}
}

func (h *Handler) Run() error {
	userInput, err := ParseUserInput()
	if err != nil {
		return err
	}

	for _, action := range h.AvailableActions {
		if action.GetInputTrigger() == userInput.ActionName {
			return action.Run(userInput.ActionParams)
		}
	}

	return fmt.Errorf("Command not implemented.")
}

func (h *Handler) DisplayHelp() {
	var entries []ui.HelpEntry
	for _, action := range h.AvailableActions {
		entries = append(entries, ui.HelpEntry{
			Name:        action.GetInputTrigger(),
			Description: action.GetDescription(),
		})
	}
	ui.DisplayHelpScreen(entries)
}

func runSync(repo *db.Repository) error {
	downloaders := []Downloader{
		NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api", repo),
		NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com", repo),
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
