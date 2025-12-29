package actions

import (
	"errors"
	"fmt"
	"math"
	actions "plutus-cli/internal/cli/actions/internal"
	"plutus-cli/internal/cli/ui"
	"strconv"
	"strings"
)

type UserDeposit struct {
	value int
}
func (d *UserDeposit) From(userInput string) error {
	userInput = strings.Replace(userInput, ",", ".", 1)

	splittedUserInput := strings.Split(userInput, ".")

	if len(splittedUserInput) < 2 {
		splittedUserInput = []string{userInput, "0"}
	}

	integerPart, fractionalPart := splittedUserInput[0], splittedUserInput[1]
	fractionalPartLength := len(fractionalPart)


	parsedIntegerPart, err := strconv.Atoi(integerPart)
	parsedFractionalPart, err := strconv.Atoi(fractionalPart)

	if err != nil {
		return err
	}
	
	d.value = parsedIntegerPart * int(math.Pow(10, float64(fractionalPartLength))) + parsedFractionalPart
	return nil
}


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
	
	deposit := UserDeposit{}
	userInput := addDepositParams[0]

	if	err := deposit.From(userInput); err != nil {
		return err
	}

	fmt.Println(deposit.value)
	return nil
}
