package testutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// for test stubbing
	osChdir     = os.Chdir
	osGetwd     = os.Getwd
	osMkdirAll  = os.MkdirAll
	osRemoveAll = os.RemoveAll
	osWriteFile = os.WriteFile
)

// AssertPaths asserts each key in files relative to base.
// Calls AssertDirPath if the key ends in "/", otherwise calls AssertFilePath.
func AssertPaths(tb testing.TB, base string, files map[string]any) {
	tb.Helper()
	if files == nil {
		return
	}
	for _, key := range sortedKeys(files) {
		value := files[key]
		path := filepath.Join(base, key)
		if isDirPath(key) {
			AssertDirPath(tb, path, value)
		} else {
			AssertFilePath(tb, path, value)
		}
	}
}

// AssertDirPath asserts the existence of a directory path.
//
//   - If value is true, then the directory SHOULD exist
//   - If value is false, then the directory SHOULD NOT exist.
func AssertDirPath(tb testing.TB, path string, value any) {
	tb.Helper()
	if exists, ok := value.(bool); ok && !exists {
		assert.NoDirExists(tb, path)
	} else {
		assert.DirExists(tb, path)
	}
}

// AssertFilePath asserts the existence of a file path.
//
//   - If value is false, then the file SHOULD NOT exist.
//   - If value is a string, then the file contents should match.
//   - If value is an int, then file permissions should match.
func AssertFilePath(tb testing.TB, path string, value any) {
	tb.Helper()
	if exists, ok := value.(bool); ok && !exists {
		// value is `false`, file _should not_ be there
		assert.NoFileExists(tb, path)
	} else {
		// file _should_ be there
		assert.FileExists(tb, path)
		if content, ok := value.(string); ok {
			// value is a string, file content should match
			buf, _ := os.ReadFile(path)
			assert.Equal(tb, content, string(buf))
		}
		if perm, ok := value.(int); ok {
			// value is an int, file permissions should match
			info, _ := os.Stat(path)
			assert.Equal(tb, perm, int(info.Mode().Perm()))
		}
	}
}

// CurrentDir delegates to os.Getwd and panics on error.
func CurrentDir() string {
	cwd, err := osGetwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

// InTempDir runs handler inside a temp dir, then returns back into the cwd.
func InTempDir(tb testing.TB, handler func(dir string)) {
	tb.Helper()

	cwd := CurrentDir()
	tmp := tb.TempDir()
	if err := osChdir(tmp); err != nil {
		panic(err)
	}
	handler(tmp)
	if err := osChdir(cwd); err != nil {
		panic(err)
	}
}

// MkdirAll delegates to os.MkdirAll and panics on error.
// The directory path will be removed when the test exits.
func MkdirAll(tb testing.TB, path string, perm fs.FileMode) {
	err := osMkdirAll(path, perm)
	if err != nil {
		panic(err)
	}
	tb.Cleanup(func() {
		RemoveAll(tb, path)
	})
}

// RemoveAll delegates to os.RemoveAll and panics on error.
func RemoveAll(tb testing.TB, path string) {
	err := osRemoveAll(path)
	if err != nil {
		panic(err)
	}
}

// WriteFile delegates to os.WriteFile and panics on error.
// The file path will be removed when the test exits.
func WriteFile(tb testing.TB, path string, data []byte, perm fs.FileMode) {
	err := osWriteFile(path, data, perm)
	if err != nil {
		panic(err)
	}
	tb.Cleanup(func() {
		RemoveAll(tb, path)
	})
}

// WritePaths creates a file for each key in files relative to base.
//
//   - The content of the file will be the map value.
//   - If the map key ends in a "/", then a directory will be created instead.
//   - Any failure to create will panic.
//   - All files and directories will be removed when the test exits.
func WritePaths(tb testing.TB, base string, files map[string]any) {
	if files == nil {
		return
	}
	// Sort so we have deterministic order in tests
	for _, key := range sortedKeys(files) {
		value := files[key]
		path := filepath.Join(base, key)
		if isDirPath(key) {
			MkdirAll(tb, path, 0755)
		} else {
			parent := filepath.Dir(path)
			if _, err := os.Stat(parent); os.IsNotExist(err) {
				MkdirAll(tb, parent, 0755)
			}
			WriteFile(tb, path, []byte(value.(string)), 0666)
		}
	}
}

// Returns whether the path refers to a directory
// (indicated by it ending in "/").
func isDirPath(path string) bool {
	return strings.HasSuffix(path, "/")
}

// Returns a sorted set of map keys.
func sortedKeys(data map[string]any) []string {
	keys := []string{}
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
