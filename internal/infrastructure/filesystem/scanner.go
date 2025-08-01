package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/priory-io/prefine/internal/domain"
)

type FileScanner struct{}

func NewFileScanner() *FileScanner {
	return &FileScanner{}
}

func (fs *FileScanner) ScanFiles(ctx context.Context, config domain.ScanConfig) ([]domain.File, error) {
	var files []domain.File

	err := filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if info.IsDir() {
			if !config.Recursive && path != config.Path {
				return filepath.SkipDir
			}
			return nil
		}

		if fs.shouldExclude(path, config.Exclude) {
			return nil
		}

		if len(config.Include) > 0 && !fs.shouldInclude(path, config.Include) {
			return nil
		}

		fileType := fs.determineFileType(info.Name())
		if fileType == "" {
			return nil
		}

		files = append(files, domain.File{
			Path: path,
			Type: fileType,
			Size: info.Size(),
		})

		return nil
	})

	return files, err
}

func (fs *FileScanner) shouldInclude(path string, patterns []string) bool {
	filename := filepath.Base(path)
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
	}
	return false
}

func (fs *FileScanner) shouldExclude(path string, patterns []string) bool {
	filename := filepath.Base(path)
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func (fs *FileScanner) determineFileType(filename string) domain.FileType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return domain.FileTypeImage
	case ".json":
		return domain.FileTypeJSON
	case ".yaml", ".yml":
		return domain.FileTypeYAML
	default:
		return ""
	}
}
