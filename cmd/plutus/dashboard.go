package main

import (
	"fmt"
	"plutus-cli/internal/db"
	"plutus-cli/internal/portfolio"
	"plutus-cli/internal/sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Displays dashboard with charts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("db").(*db.Repository)

		sync.RunSync(repo)

		data, _ := portfolio.GetChartData(repo)
		fmt.Println(data)

		return nil
	},
}
