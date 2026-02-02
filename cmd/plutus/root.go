package main

import (
	"context"
	"fmt"
	"os"
	"plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"

	"github.com/spf13/cobra"
)

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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ui.CheckForUpdates(Version)
	},
}

func Execute(database *db.Repository) {
	ctx := context.WithValue(context.Background(), "db", database)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
