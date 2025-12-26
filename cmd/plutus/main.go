package main

import (
	"flag"
	"plutus-cli/internal/cli/ui"
)

func main() {
	ui.PrintBanner()

	flag.Bool("help", false, "help flag")

	flag.Parse()

	ui.PrintHelp()
}
