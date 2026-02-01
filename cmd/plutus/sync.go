package main

import (
	"fmt"
	"plutus-cli/internal/cli/actions"
	"plutus-cli/internal/db"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches up-to-date market data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSync(repo)
	},
}

func runSync(repo *db.Repository) error {
	downloaders := []actions.Downloader{
		actions.NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api", repo),
		actions.NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com", repo),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(downloaders))

	for _, downloader := range downloaders {
		wg.Add(1)
		go func(d actions.Downloader) {
			defer wg.Done()
			if err := d.SyncData(); err != nil {
				errChan <- fmt.Errorf("%s failed: %w", d.GetName(), err)
			}
		}(downloader)
	}

	wg.Wait()
	close(errChan)

	var combinedError error
	for err := range errChan {
		if combinedError == nil {
			combinedError = err
		} else {
			fmt.Println("Error:", err)
		}
	}

	if combinedError != nil {
		return combinedError
	}

	fmt.Println("Sync completed successfully.")
	return nil
}
