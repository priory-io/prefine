package processors

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/priory-io/prefine/internal/domain"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) CanProcess(file domain.File) bool {
	return file.Type == domain.FileTypeImage
}

func (p *ImageProcessor) Process(
	ctx context.Context,
	file domain.File,
	config domain.OptimizeConfig,
) domain.OptimizeResult {
	start := time.Now()

	result := domain.OptimizeResult{
		File:         file,
		OriginalSize: file.Size,
		Duration:     0,
	}

	if config.DryRun {
		result.NewSize = result.OriginalSize
		result.Duration = time.Since(start)
		return result
	}

	img, err := imaging.Open(file.Path)
	if err != nil {
		result.Error = fmt.Errorf("failed to open image: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	optimized := p.resizeIfNeeded(img, config)

	tempPath := file.Path + ".tmp"
	err = p.saveOptimizedImage(optimized, tempPath, file.Path, config)
	if err != nil {
		result.Error = fmt.Errorf("failed to save optimized image: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	tempStat, err := os.Stat(tempPath)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("failed to stat optimized file: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	result.NewSize = tempStat.Size()
	result.Savings = result.OriginalSize - result.NewSize

	if result.NewSize >= result.OriginalSize {
		os.Remove(tempPath)
		result.NewSize = result.OriginalSize
		result.Savings = 0
		result.Duration = time.Since(start)
		return result
	}

	err = os.Rename(tempPath, file.Path)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("failed to replace original file: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Duration = time.Since(start)
	return result
}

func (p *ImageProcessor) resizeIfNeeded(img image.Image, config domain.OptimizeConfig) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > config.MaxWidth || height > config.MaxHeight {
		return imaging.Fit(img, config.MaxWidth, config.MaxHeight, imaging.Lanczos)
	}

	return img
}

func (p *ImageProcessor) saveOptimizedImage(img image.Image, outputPath, originalPath string, config domain.OptimizeConfig) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(originalPath))
	switch ext {
	case ".jpg", ".jpeg":
		opts := &jpeg.Options{Quality: config.Quality}
		return jpeg.Encode(out, img, opts)
	case ".png":
		encoder := &png.Encoder{CompressionLevel: png.BestCompression}
		return encoder.Encode(out, img)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
}
