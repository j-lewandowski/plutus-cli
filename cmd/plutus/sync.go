package main

import (
	"plutus-cli/internal/db"
	"plutus-cli/internal/sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches up-to-date market data",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("db").(*db.Repository)
		return sync.RunSync(repo)
	},
}
