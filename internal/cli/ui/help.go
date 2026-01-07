package ui

import "fmt"

func PrintHelp() {
	fmt.Print("Example usage: \t plutus [options] COMMAND \n\n")
	fmt.Println("Avaliable Commands:")
	fmt.Println("\t - add \t\t Allows user to add deposit event.")
	fmt.Print("Global Options: \n")
	fmt.Println("--help \t See more information on a command")
}

func DisplayHelpScreen() {
	PrintBanner()
	PrintHelp()
}
