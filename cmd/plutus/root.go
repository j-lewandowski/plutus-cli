package main

import (
	"fmt"
	"os"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"

	"github.com/spf13/cobra"
)

var repo *db.Repository

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()

		fmt.Println(cmd.UsageString())
	})
}

var rootCmd = &cobra.Command{
	Use:   "plutus",
	Short: "Plutus CLI - personal investment tracker",
	Long:  `Plutus is a local-first CLI tool written in Go for tracking long-term investments and visualizing portfolio performance over time.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(database *db.Repository) {
	repo = database
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
