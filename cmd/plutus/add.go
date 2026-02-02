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
	Short: "Adds a new deposit",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := db.NewDepositParams{
			DepositAmount: args[0],
		}
		if len(args) == 2 {
			params.DepositDate = args[1]
		}

		deposit := db.UserDeposit{}
		if err := deposit.From(params); err != nil {
			return err
		}

		if err := repo.AddDeposit(deposit); err != nil {
			return err
		}

		fmt.Println("Deposit added!")
		return nil
	},
}
