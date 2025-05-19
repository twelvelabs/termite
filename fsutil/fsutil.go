package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	// DefaultDirMode grants `rwx------`.
	DefaultDirMode = 0700
	// DefaultFileMode grants `rw-------`.
	DefaultFileMode = 0600
)

var (
	osUserHomeDir = os.UserHomeDir
	filepathAbs   = filepath.Abs
	osMkdirAll    = os.MkdirAll
	osWriteFile   = os.WriteFile
)

func NoPathExists(path string) bool {
	_, err := os.Stat(path)
	// for some reason os.ErrInvalid sometimes != syscall.EINVAL :shrug:
	if errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, os.ErrInvalid) ||
		errors.Is(err, syscall.EINVAL) {
		return true
	}
	return false
}

func PathExists(path string) bool {
	return !NoPathExists(path)
}

// NormalizePath ensures that name is an absolute path.
// Environment variables (and the ~ string) are expanded.
func NormalizePath(name string) (string, error) {
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", nil
	}

	// Replace ENV vars
	normalized = os.ExpandEnv(normalized)

	// Replace ~
	if strings.HasPrefix(normalized, "~") {
		home, err := osUserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to normalize %s: %w", name, err)
		}
		normalized = home + strings.TrimPrefix(normalized, "~")
	}

	// Ensure abs path
	normalized, err := filepathAbs(normalized)
	if err != nil {
		return "", fmt.Errorf("unable to normalize %s: %w", name, err)
	}

	return normalized, nil
}

// EnsureDirWritable ensures that path is a writable directory.
// Will attempt to create a new directory if path does not exist.
func EnsureDirWritable(path string) error {
	// Ensure dir exists (and IsDir).
	err := osMkdirAll(path, DefaultDirMode)
	if err != nil {
		return fmt.Errorf("ensure dir: %w", err)
	}

	f := filepath.Join(path, ".touch")
	if err := osWriteFile(f, []byte(""), DefaultFileMode); err != nil {
		return fmt.Errorf("ensure writable: %w", err)
	}
	defer func() {
		_ = os.Remove(f)
	}()

	return nil
}
