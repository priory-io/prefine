package cmd

import (
	"context"
	"os"

	"github.com/priory-io/prefine/internal/application"
	"github.com/priory-io/prefine/internal/domain"
	"github.com/priory-io/prefine/internal/infrastructure/filesystem"
	"github.com/priory-io/prefine/internal/infrastructure/logging"
	"github.com/priory-io/prefine/internal/infrastructure/processors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prefine [path]",
	Short: "Optimize images, JSON, YAML and other web development files",
	Long: `prefine is a CLI tool designed to optimize various file types commonly
found in web development projects including images (PNG, JPG, WebP),
configuration files (JSON, YAML), and more.

If no path is provided, optimizes the current directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		recursive, _ := cmd.Flags().GetBool("recursive")
		quality, _ := cmd.Flags().GetInt("quality")
		maxWidth, _ := cmd.Flags().GetInt("max-width")
		maxHeight, _ := cmd.Flags().GetInt("max-height")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			logger := logging.NewConsoleLogger(verbose)
			logger.Error("Path '%s' does not exist", path)
			os.Exit(1)
		}

		scanner := filesystem.NewFileScanner()
		processors := []domain.FileProcessor{
			processors.NewImageProcessor(),
		}
		logger := logging.NewConsoleLogger(verbose)

		service := application.NewOptimizationService(scanner, processors, logger)

		scanConfig := domain.ScanConfig{
			Path:      path,
			Include:   include,
			Exclude:   exclude,
			Recursive: recursive,
		}

		optimizeConfig := domain.OptimizeConfig{
			Quality:   quality,
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
			Verbose:   verbose,
			DryRun:    dryRun,
		}

		ctx := context.Background()
		report := service.OptimizeFiles(ctx, scanConfig, optimizeConfig)

		if verbose {
			for _, result := range report.Results {
				if result.Error != nil {
					logger.Error("Failed to optimize %s: %v", result.File.Path, result.Error)
				} else if result.WasOptimized() {
					logger.Info("Optimized %s: %d bytes (%.1f%%) saved in %v",
						result.File.Path, result.Savings, result.SavingsPercent(), result.Duration)
				}
			}
		}

		if report.FailedFiles > 0 {
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Show what would be optimized without making changes")

	rootCmd.Flags().StringSliceP("include", "i", []string{}, "File patterns to include (e.g., *.png,*.jpg)")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "File patterns to exclude")
	rootCmd.Flags().BoolP("recursive", "r", true, "Recursively process subdirectories")
	rootCmd.Flags().IntP("quality", "q", 85, "JPEG quality (1-100)")
	rootCmd.Flags().Int("max-width", 1920, "Maximum width for image resizing")
	rootCmd.Flags().Int("max-height", 1080, "Maximum height for image resizing")
}
