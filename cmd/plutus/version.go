package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print detailed version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Plutus CLI:  %s\n", Version)
		fmt.Printf("Git Commit:  %s\n", Commit)
		fmt.Printf("Build Time:  %s\n", BuildTime)
		fmt.Printf("Go version:  %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
