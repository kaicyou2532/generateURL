package local

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Saver stores files on the local filesystem.
type Saver struct {
	Dir string
}

// New creates a Saver ensuring the target directory exists.
func New(dir string) (*Saver, error) {
	if dir == "" {
		return nil, fmt.Errorf("upload directory is required")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	return &Saver{Dir: dir}, nil
}

// Save writes the given file to disk and returns the public file name.
func (s *Saver) Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	if file == nil || header == nil {
		return "", fmt.Errorf("file and header are required")
	}
	defer file.Close()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	ext := filepath.Ext(header.Filename)
	name := strings.ToLower(uuid.NewString()) + ext
	fullPath := filepath.Join(s.Dir, name)

	tmpFile, err := os.CreateTemp(s.Dir, name+".*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("write file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("close file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), fullPath); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("finalize file: %w", err)
	}

	return name, nil
}
