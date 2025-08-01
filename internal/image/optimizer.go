package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type OptimizeOptions struct {
	Quality     int
	MaxWidth    int
	MaxHeight   int
	Progressive bool
	Verbose     bool
	DryRun      bool
}

type OptimizeResult struct {
	OriginalPath string
	OriginalSize int64
	NewSize      int64
	Savings      int64
	Error        error
}

func (r OptimizeResult) SavingsPercent() float64 {
	if r.OriginalSize == 0 {
		return 0
	}
	return float64(r.Savings) / float64(r.OriginalSize) * 100
}

var DefaultOptions = OptimizeOptions{
	Quality:     85,
	MaxWidth:    1920,
	MaxHeight:   1080,
	Progressive: true,
	Verbose:     false,
	DryRun:      false,
}

func OptimizeImage(imagePath string, options OptimizeOptions) OptimizeResult {
	result := OptimizeResult{
		OriginalPath: imagePath,
	}

	stat, err := os.Stat(imagePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat file: %w", err)
		return result
	}
	result.OriginalSize = stat.Size()

	if options.DryRun {
		result.NewSize = result.OriginalSize
		return result
	}

	img, err := imaging.Open(imagePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to open image: %w", err)
		return result
	}

	optimized := img
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > options.MaxWidth || height > options.MaxHeight {
		optimized = imaging.Fit(img, options.MaxWidth, options.MaxHeight, imaging.Lanczos)
		if options.Verbose {
			fmt.Printf("Resized %s from %dx%d to %dx%d\n",
				filepath.Base(imagePath), width, height,
				optimized.Bounds().Dx(), optimized.Bounds().Dy())
		}
	}

	tempPath := imagePath + ".tmp"
	err = saveOptimizedImage(optimized, tempPath, imagePath, options)
	if err != nil {
		result.Error = fmt.Errorf("failed to save optimized image: %w", err)
		return result
	}

	tempStat, err := os.Stat(tempPath)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("failed to stat optimized file: %w", err)
		return result
	}

	result.NewSize = tempStat.Size()
	result.Savings = result.OriginalSize - result.NewSize

	if result.NewSize >= result.OriginalSize {
		os.Remove(tempPath)
		result.NewSize = result.OriginalSize
		result.Savings = 0
		if options.Verbose {
			fmt.Printf("No improvement for %s, keeping original\n", filepath.Base(imagePath))
		}
		return result
	}

	err = os.Rename(tempPath, imagePath)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("failed to replace original file: %w", err)
		return result
	}

	return result
}

func saveOptimizedImage(img image.Image, outputPath, originalPath string, options OptimizeOptions) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(originalPath))
	switch ext {
	case ".jpg", ".jpeg":
		opts := &jpeg.Options{Quality: options.Quality}
		return jpeg.Encode(out, img, opts)
	case ".png":
		encoder := &png.Encoder{CompressionLevel: png.BestCompression}
		return encoder.Encode(out, img)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
}

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}
