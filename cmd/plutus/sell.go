package main

import (
	"fmt"
	"plutus-cli/internal/db"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sellCmd)
}

var sellCmd = &cobra.Command{
	Use:   "sell <amount> [date]",
	Short: "Records a sell / withdrawal",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("db").(*db.Repository)

		params := db.NewSellParams{
			SellAmount: args[0],
		}
		if len(args) == 2 {
			params.SellDate = args[1]
		}

		sell := db.UserSell{}
		if err := sell.From(params); err != nil {
			return err
		}

		if err := repo.AddSell(sell); err != nil {
			return err
		}

		fmt.Println("Sell recorded!")
		return nil
	},
}
