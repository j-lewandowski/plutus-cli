package main

import (
	"flag"
	"fmt"
	"plutus-cli/internal/cli/actions"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
)

func main() {
	repo, err := db.InitDb()

	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	helpFlag := flag.Bool("help", false, "help flag")

	flag.Parse()

	if *helpFlag {
		ui.DisplayHelpScreen()
		return
	}

	err = actions.HandleUserAction(repo)

	defer repo.Close()

	if err != nil {
		fmt.Println("Couln't perform this operation because of an error:", err)
		return
	}
}
