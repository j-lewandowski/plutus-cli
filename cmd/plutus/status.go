package main

import (
	"fmt"
	"plutus-cli/internal/db" // 1. Add this import
	"plutus-cli/internal/portfolio"
	"plutus-cli/internal/sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays current portfolio value",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("db").(*db.Repository)

		total, err := repo.GetOverallDepositInEurocents()
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}

		if total == 0 {
			fmt.Println("Your portfolio is empty. Add your first deposit using: plutus add <amount>")
			return nil
		}

		_ = sync.RunSync(repo)

		report, err := portfolio.CalculatePortfolio(repo)
		if err != nil {
			return fmt.Errorf("could not calculate portfolio: %w", err)
		}

		totalInvested := float64(report.TotalInvestedInEurocents) / 100.0
		currentValue := float64(report.CurrentValueInEurocents) / 100.0
		fmt.Printf("Total Invested: %.2f EUR\n", totalInvested)
		fmt.Printf("Current Value:  %.2f EUR\n", currentValue)
		if totalInvested > 0 {
			profitLossPct := (currentValue - totalInvested) / totalInvested * 100.0
			fmt.Printf("Profit/Loss:    %.2f%%\n", profitLossPct)
		}
		return nil
	},
}
