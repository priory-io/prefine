package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	fmt.Println("prefine v0.1.0")
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
