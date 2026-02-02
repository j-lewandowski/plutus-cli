package main

import (
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
		return sync.RunSync(repo)
	},
}
