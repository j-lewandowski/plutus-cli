package ui

import "fmt"

func PrintHelp() {
	fmt.Print("Example usage: \t plutus [options] COMMAND \n\n")
	fmt.Println("Available Commands:")
	fmt.Println("\t - add \t\t Allows user to add deposit event.")
	fmt.Println("\t - sync \t Syncs the CLI with up-to-date market data.")
	fmt.Println("\t - status \t Displays current portfolio value and profit/loss percentage.")
	fmt.Print("Global Options: \n")
	fmt.Println("--help \t See more information on a command")
}

func DisplayHelpScreen() {
	PrintBanner()
	PrintHelp()
}
