package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
)

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print the version number of prefine",
	Long:   "Print the version number of prefine",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("prefine %s\n", Version)
	fmt.Printf("Built: %s\n", BuildTime)
	fmt.Printf("Commit: %s\n", GitCommit)
	fmt.Printf("Go: %s\n", GoVersion)
	os.Exit(0)
}

func printVersionShort() {
	fmt.Printf("prefine %s\n", Version)
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
