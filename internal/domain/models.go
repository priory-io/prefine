package domain

import "time"

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeJSON  FileType = "json"
	FileTypeYAML  FileType = "yaml"
)

type File struct {
	Path string
	Type FileType
	Size int64
}

type OptimizeConfig struct {
	Quality   int
	MaxWidth  int
	MaxHeight int
	DryRun    bool
	Verbose   bool
}

type OptimizeResult struct {
	File         File
	OriginalSize int64
	NewSize      int64
	Savings      int64
	Duration     time.Duration
	Error        error
}

func (r OptimizeResult) SavingsPercent() float64 {
	if r.OriginalSize == 0 {
		return 0
	}
	return float64(r.Savings) / float64(r.OriginalSize) * 100
}

func (r OptimizeResult) WasOptimized() bool {
	return r.Savings > 0 && r.Error == nil
}

type ScanConfig struct {
	Path      string
	Include   []string
	Exclude   []string
	Recursive bool
}

type OptimizationReport struct {
	TotalFiles     int
	OptimizedFiles int
	FailedFiles    int
	TotalSavings   int64
	Duration       time.Duration
	Results        []OptimizeResult
}

func (r OptimizationReport) TotalSavingsPercent() float64 {
	var totalOriginal int64
	for _, result := range r.Results {
		totalOriginal += result.OriginalSize
	}
	if totalOriginal == 0 {
		return 0
	}
	return float64(r.TotalSavings) / float64(totalOriginal) * 100
}
