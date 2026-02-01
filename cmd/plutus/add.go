package main

import (
	"fmt"
	"plutus-cli/internal/db"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add <amount> [date]",
	Short: "Adds a new deposit to your portfolio",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Syncing market data before adding deposit...")
		if err := runSync(repo); err != nil {
			return fmt.Errorf("could not sync data before adding: %w", err)
		}

		amount := args[0]
		date := ""
		if len(args) > 1 {
			date = args[1]
		}

		params := db.NewDepositParams{
			DepositAmount: amount,
			DepositDate:   date,
		}

		deposit := db.UserDeposit{}
		if err := deposit.From(params); err != nil {
			return err
		}

		if err := repo.AddDeposit(deposit); err != nil {
			return err
		}

		fmt.Println("Deposit added successfully!")
		return nil
	},
}
