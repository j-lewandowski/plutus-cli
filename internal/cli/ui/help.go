package ui

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type HelpEntry struct {
	Name        string
	Description string
}

func PrintHelp(commands []HelpEntry) {
	fmt.Print("Example usage: \t plutus [options] COMMAND \n\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Available Commands:")
	for _, cmd := range commands {
		fmt.Fprintf(w, "\t- %s\t%s\n", cmd.Name, cmd.Description)
	}
	w.Flush()

	fmt.Print("\nGlobal Options: \n")
	fmt.Println("--help \t See more information on a command")
}

func DisplayHelpScreen(commands []HelpEntry) {
	PrintBanner()
	PrintHelp(commands)
}
