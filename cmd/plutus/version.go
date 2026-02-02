var (
	Version   = "dev"
	Commit    = ""
	BuildTime = "unknown"
)

func Execute(r *db.Repository) {
	repo = r
	rootCmd.Version = Version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

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