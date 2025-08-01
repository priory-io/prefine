package domain

import "errors"

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileNotFound        = errors.New("file not found")
	ErrInvalidConfig       = errors.New("invalid configuration")
)
