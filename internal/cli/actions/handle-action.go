package actions

import (
	"errors"
	"fmt"
	actions "plutus-cli/internal/cli/actions/internal"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db/models"
	"plutus-cli/internal/db/queries"
)

func HandleUserAction() error {

	userInput, err := actions.ParseUserInput()

	if err != nil {
		return err
	}

	switch userInput.ActionName {
	case "add":
		if err := handleAddDeposit(userInput.ActionParams); err != nil {
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

	parsedDepositParams := models.AddDepositParams{
		DepositAmount: addDepositParams[0],
	}

	if len(addDepositParams) == 2 {
		parsedDepositParams.DepositDate = addDepositParams[1]
	}

	deposit := models.UserDeposit{}

	if err := deposit.From(parsedDepositParams); err != nil {
		return err
	}

	if err := queries.AddDeposit(deposit); err != nil {
		return err
	}

	fmt.Println("Deposit added!")

	return nil
}
