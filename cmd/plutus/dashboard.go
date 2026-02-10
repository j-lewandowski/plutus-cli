package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	cli "plutus-cli/internal/cli/ui"
	"plutus-cli/internal/db"
	"plutus-cli/internal/portfolio"
	"plutus-cli/internal/sync"
	"plutus-cli/internal/ui"
	"runtime"

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

		http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
			points, err := portfolio.GetChartData(repo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(points)
		})

		assets, err := ui.GetFileSystem()
		if err != nil {
			return fmt.Errorf("Could not load UI assets: %w", err)
		}

		http.Handle("/", http.FileServer(http.FS(assets)))

		cli.PrintBanner()

		port := ":8055"
		url := "http://localhost" + port

		fmt.Printf("🚀 Dashboard starting at %s\n", url)
		fmt.Println("Press Ctrl+C to close.")

		openBrowser(url)

		if err := http.ListenAndServe(port, nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}

		return nil
	},
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Printf("Nie udało się otworzyć przeglądarki: %v\n", err)
	}
}
