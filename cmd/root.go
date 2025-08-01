package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prefine",
	Short: "Optimize images, JSON, YAML and other web development files",
	Long: `prefine is a CLI tool designed to optimize various file types commonly
found in web development projects including images (PNG, JPG, WebP),
configuration files (JSON, YAML), and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Show what would be optimized without making changes")
}
