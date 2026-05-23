package utils

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidRange     = errors.New("invalid line range")
	ErrInvalidPatch     = errors.New("invalid patch format")
	ErrInvalidPath      = errors.New("invalid path")
)
