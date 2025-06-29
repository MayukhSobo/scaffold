package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveLogFilePath creates the full path for the log file.
// If the directory is not absolute, it's considered relative to the current working directory.
func ResolveLogFilePath(directory, filename string) string {
	if directory == "" {
		directory = "logs" // Default directory
	}

	if !filepath.IsAbs(directory) {
		cwd, err := os.Getwd()
		if err != nil {
			// On error, fallback to a simple logs directory
			return filepath.Join("logs", filename)
		}
		directory = filepath.Join(cwd, directory)
	}

	return filepath.Join(directory, filename)
}

// EnsureLogDirectory creates the log directory if it doesn't exist.
func EnsureLogDirectory(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}
	return nil
}
