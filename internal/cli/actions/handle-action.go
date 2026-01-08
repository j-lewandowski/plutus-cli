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
		if err := handleAddDeposit(userInput.ActionParams); err != nil {
			return err
		}
		return nil

	case "sync":
		if err := handleSync(); err != nil {
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
	downloaders := []Downloader{NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api")}

	for _, downloader := range downloaders {
		err := downloader.SyncData()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
