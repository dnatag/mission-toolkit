package utils

import (
	"fmt"

	"github.com/spf13/afero"
)

// CopyFile copies a file from src to dst using the provided filesystem.
// It returns an error if reading the source file or writing to the destination fails.
func CopyFile(fs afero.Fs, src, dst string) error {
	content, err := afero.ReadFile(fs, src)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	if err := afero.WriteFile(fs, dst, content, 0644); err != nil {
		return fmt.Errorf("writing destination file: %w", err)
	}

	return nil
}
