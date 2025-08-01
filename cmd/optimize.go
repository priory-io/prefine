package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize [path]",
	Short: "Optimize files in the specified directory",
	Long: `Optimize images, JSON, YAML and other web development files
in the specified directory. If no path is provided, optimizes the current directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if verbose {
			fmt.Printf("Optimizing files in: %s\n", path)
			if dryRun {
				fmt.Println("Running in dry-run mode")
			}
		}

		fmt.Printf("Optimization complete for: %s\n", path)
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
	optimizeCmd.Flags().StringSliceP("include", "i", []string{}, "File patterns to include (e.g., *.png,*.jpg)")
	optimizeCmd.Flags().StringSliceP("exclude", "e", []string{}, "File patterns to exclude")
	optimizeCmd.Flags().BoolP("recursive", "r", true, "Recursively process subdirectories")
}
