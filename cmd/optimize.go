package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/priory-io/prefine/internal/image"
	"github.com/priory-io/prefine/internal/scanner"
	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize [path]",
	Short: "Optimize files in the specified directory",
	Long: `Optimize images and other files
in the specified directory. If no path is provided, optimizes the current directory.`,
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

		if verbose {
			fmt.Printf("Optimizing files in: %s\n", path)
			if dryRun {
				fmt.Println("Running in dry-run mode")
			}
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Error: Path '%s' does not exist\n", path)
			os.Exit(1)
		}

		fileScanner := scanner.NewFileScanner(include, exclude, recursive)
		imageFiles, err := fileScanner.ScanForImages(path)
		if err != nil {
			fmt.Printf("Error scanning for images: %v\n", err)
			os.Exit(1)
		}

		if len(imageFiles) == 0 {
			fmt.Println("No image files found to optimize")
			return
		}

		fmt.Printf("Found %d image files to optimize\n", len(imageFiles))

		options := image.OptimizeOptions{
			Quality:   quality,
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
			Verbose:   verbose,
			DryRun:    dryRun,
		}

		start := time.Now()
		var totalSavings int64
		var optimizedCount int
		var errorCount int

		for _, imagePath := range imageFiles {
			if verbose {
				fmt.Printf("Processing: %s\n", imagePath)
			}

			result := image.OptimizeImage(imagePath, options)
			if result.Error != nil {
				fmt.Printf("Error optimizing %s: %v\n", imagePath, result.Error)
				errorCount++
				continue
			}

			if result.Savings > 0 {
				optimizedCount++
				totalSavings += result.Savings
				if verbose || !dryRun {
					fmt.Printf("Optimized %s: %d bytes (%.1f%%) saved\n",
						imagePath, result.Savings, result.SavingsPercent())
				}
			}
		}

		duration := time.Since(start)
		fmt.Printf("\nOptimization complete in %v\n", duration)
		fmt.Printf("Files processed: %d\n", len(imageFiles))
		fmt.Printf("Files optimized: %d\n", optimizedCount)
		if errorCount > 0 {
			fmt.Printf("Errors encountered: %d\n", errorCount)
		}
		if totalSavings > 0 {
			fmt.Printf("Total space saved: %d bytes\n", totalSavings)
		}
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
	optimizeCmd.Flags().StringSliceP("include", "i", []string{}, "File patterns to include (e.g., *.png,*.jpg)")
	optimizeCmd.Flags().StringSliceP("exclude", "e", []string{}, "File patterns to exclude")
	optimizeCmd.Flags().BoolP("recursive", "r", true, "Recursively process subdirectories")
	optimizeCmd.Flags().IntP("quality", "q", 85, "JPEG quality (1-100)")
	optimizeCmd.Flags().Int("max-width", 1920, "Maximum width for image resizing")
	optimizeCmd.Flags().Int("max-height", 1080, "Maximum height for image resizing")
}
