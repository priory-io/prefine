package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/priory-io/prefine/internal/image"
)

type FileScanner struct {
	Include   []string
	Exclude   []string
	Recursive bool
}

func NewFileScanner(include, exclude []string, recursive bool) *FileScanner {
	return &FileScanner{
		Include:   include,
		Exclude:   exclude,
		Recursive: recursive,
	}
}

func (fs *FileScanner) ScanForImages(rootPath string) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !fs.Recursive && path != rootPath {
				return filepath.SkipDir
			}
			return nil
		}

		if !image.IsImageFile(info.Name()) {
			return nil
		}

		if fs.shouldExclude(path) {
			return nil
		}

		if len(fs.Include) > 0 && !fs.shouldInclude(path) {
			return nil
		}

		imageFiles = append(imageFiles, path)
		return nil
	})

	return imageFiles, err
}

func (fs *FileScanner) shouldInclude(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range fs.Include {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
	}
	return false
}

func (fs *FileScanner) shouldExclude(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range fs.Exclude {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}
