package main

import (
	"flag"
	"fmt"
	"plutus-cli/internal/cli/actions"
	"plutus-cli/internal/db"
)

func main() {
	repo, err := db.InitDb()

	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer repo.Close()

	helpFlag := flag.Bool("help", false, "help flag")

	flag.Parse()

	handler := actions.NewHandler(repo)

	if *helpFlag {
		handler.DisplayHelp()
		return
	}

	err = handler.Run()

	if err != nil {
		fmt.Println("Couln't perform this operation because of an error:", err)
		return
	}
}
