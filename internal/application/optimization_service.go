package application

import (
	"context"
	"sync"
	"time"

	"github.com/priory-io/prefine/internal/domain"
)

type OptimizationService struct {
	scanner    domain.FileScanner
	processors []domain.FileProcessor
	logger     domain.Logger
}

func NewOptimizationService(
	scanner domain.FileScanner,
	processors []domain.FileProcessor,
	logger domain.Logger,
) *OptimizationService {
	return &OptimizationService{
		scanner:    scanner,
		processors: processors,
		logger:     logger,
	}
}

func (s *OptimizationService) OptimizeFiles(
	ctx context.Context,
	scanConfig domain.ScanConfig,
	optimizeConfig domain.OptimizeConfig,
) domain.OptimizationReport {
	start := time.Now()

	s.logger.Info("Starting file optimization in: %s", scanConfig.Path)

	files, err := s.scanner.ScanFiles(ctx, scanConfig)
	if err != nil {
		s.logger.Error("Failed to scan files: %v", err)
		return domain.OptimizationReport{
			Duration: time.Since(start),
		}
	}

	if len(files) == 0 {
		s.logger.Info("No files found to optimize")
		return domain.OptimizationReport{
			Duration: time.Since(start),
		}
	}

	s.logger.Info("Found %d files to process", len(files))

	results := s.processFiles(ctx, files, optimizeConfig)

	report := s.buildReport(results, time.Since(start))
	s.logSummary(report)

	return report
}

func (s *OptimizationService) processFiles(
	ctx context.Context,
	files []domain.File,
	config domain.OptimizeConfig,
) []domain.OptimizeResult {
	results := make([]domain.OptimizeResult, len(files))
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(index int, f domain.File) {
			defer wg.Done()
			results[index] = s.processFile(ctx, f, config)
		}(i, file)
	}

	wg.Wait()
	return results
}

func (s *OptimizationService) processFile(
	ctx context.Context,
	file domain.File,
	config domain.OptimizeConfig,
) domain.OptimizeResult {
	processor := s.findProcessor(file)
	if processor == nil {
		return domain.OptimizeResult{
			File:  file,
			Error: domain.ErrUnsupportedFileType,
		}
	}

	if config.Verbose {
		s.logger.Info("Processing: %s", file.Path)
	}

	return processor.Process(ctx, file, config)
}

func (s *OptimizationService) findProcessor(file domain.File) domain.FileProcessor {
	for _, processor := range s.processors {
		if processor.CanProcess(file) {
			return processor
		}
	}
	return nil
}

func (s *OptimizationService) buildReport(
	results []domain.OptimizeResult,
	duration time.Duration,
) domain.OptimizationReport {
	var optimizedFiles, failedFiles int
	var totalSavings int64

	for _, result := range results {
		if result.Error != nil {
			failedFiles++
		} else if result.WasOptimized() {
			optimizedFiles++
		}
		totalSavings += result.Savings
	}

	return domain.OptimizationReport{
		TotalFiles:     len(results),
		OptimizedFiles: optimizedFiles,
		FailedFiles:    failedFiles,
		TotalSavings:   totalSavings,
		Duration:       duration,
		Results:        results,
	}
}

func (s *OptimizationService) logSummary(report domain.OptimizationReport) {
	s.logger.Info("Optimization complete in %v", report.Duration)
	s.logger.Info("Files processed: %d", report.TotalFiles)
	s.logger.Info("Files optimized: %d", report.OptimizedFiles)

	if report.FailedFiles > 0 {
		s.logger.Info("Errors encountered: %d", report.FailedFiles)
	}

	if report.TotalSavings > 0 {
		s.logger.Info("Total space saved: %d bytes (%.1f%%)",
			report.TotalSavings, report.TotalSavingsPercent())
	}
}
