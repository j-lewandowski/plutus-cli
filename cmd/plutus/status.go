package main

import (
	"fmt"
	"plutus-cli/internal/portfolio"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays current portfolio value",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runSync(repo); err != nil {
			return err
		}

		report, err := portfolio.CalculatePortfolio(repo)
		if err != nil {
			return err
		}

		for _, warning := range report.Warnings {
			fmt.Println(warning)
		}

		fmt.Printf("Total Invested: %.2f EUR\n", float64(report.TotalInvestedInEurocents)/100.0)
		fmt.Printf("Current Value:  %.2f EUR\n", float64(report.CurrentValueInEurocents)/100.0)
		fmt.Printf("Profit/Loss:    %.2f EUR (%.2f%%)\n", float64(report.ProfitValueInEurocents)/100.0, report.ProfitPercent)

		if report.HasExchangeRate {
			fmt.Println("---------------------------")
			fmt.Printf("Rate (1 EUR):   %.4f PLN\n", float64(report.RateEURtoPLNInGrosz)/100.0)
			fmt.Printf("Assets Value:   %.2f PLN\n", float64(report.CurrentValueInGrosz)/100.0)
			fmt.Printf("Profit/Loss: %.2f PLN\n", float64(report.ProfitValueInGrosz)/100.0)
		}

		return nil
	},
}
