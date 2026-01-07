package actions

import (
	"errors"
	"os"
	"slices"
	"strings"
)

// @TODO - figure out better way to handle enums because go doesn't really have them :c
var AvaliableCommands = []string{"add", "sync"}

type UserInput struct {
	ActionName   string
	ActionParams []string
}

func ParseUserInput() (UserInput, error) {
	args := os.Args

	if len(args) < 3 {
		return UserInput{
			ActionName: "help",
		}, nil
	}

	params := args[2:]

	command, err := ValidateCommand(args[1])

	if err != nil {
		return UserInput{}, err
	}

	return UserInput{
		ActionName:   command,
		ActionParams: params,
	}, nil
}

func ValidateCommand(command string) (string, error) {
	lowercaseCommand := strings.ToLower(command)

	if !slices.Contains(AvaliableCommands, lowercaseCommand) {
		return "", errors.New("Command not implemented.")
	}

	return lowercaseCommand, nil
}
