package main

import (
	"flag"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
)

func main() {
	_, err := db.InitDb()

	if err != nil {
		println("Can't initialize database:", err)
		return 
	}

	flag.Bool("help", false, "help flag")

	flag.Parse()

	ui.PrintBanner()
	ui.PrintHelp()
}
