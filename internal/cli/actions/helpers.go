package actions

import (
	"os"
	"strings"
	"time"
)

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
	command := strings.ToLower(args[1])

	return UserInput{
		ActionName:   command,
		ActionParams: params,
	}, nil
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
