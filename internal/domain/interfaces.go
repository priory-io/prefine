package domain

import "context"

type FileScanner interface {
	ScanFiles(ctx context.Context, config ScanConfig) ([]File, error)
}

type FileProcessor interface {
	CanProcess(file File) bool
	Process(ctx context.Context, file File, config OptimizeConfig) OptimizeResult
}

type OptimizationService interface {
	OptimizeFiles(ctx context.Context, scanConfig ScanConfig, optimizeConfig OptimizeConfig) OptimizationReport
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
