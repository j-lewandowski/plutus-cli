package actions

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

// @TODO - figure out better way to handle enums because go doesn't really have them :c
var AvailableCommands = []string{"add", "sync", "status"}

type UserInput struct {
	ActionName   string
	ActionParams []string
}

func ParseUserInput() (UserInput, error) {
	args := os.Args

	if len(args) < 2 {
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

	if !slices.Contains(AvailableCommands, lowercaseCommand) {
		return "", fmt.Errorf("Command not implemented.")
	}

	return lowercaseCommand, nil
}

func DaysUntilToday(startDate time.Time) []time.Time {
	var days []time.Time

	now := time.Now()

	cursor := startDate

	for !cursor.After(now) {
		days = append(days, cursor)
		cursor = cursor.AddDate(0, 0, 1)
	}

	return days
}
