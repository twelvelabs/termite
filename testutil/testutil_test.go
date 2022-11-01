package testutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func TestAssertFiles(t *testing.T) {
	InTempDir(t, func(dir string) {
		WritePaths(t, dir, map[string]any{
			"aaa/":       true,
			"bbb/ccc/":   true,
			"bin/aaa.sh": "aaa",
			"bin/bbb.sh": "bbb",
			"hello.txt":  "hello",
		})

		_ = os.Chmod(filepath.Join(dir, "bin", "aaa.sh"), 0600)
		_ = os.Chmod(filepath.Join(dir, "bin", "bbb.sh"), 0600)

		AssertPaths(t, dir, map[string]any{
			"aaa/":        true,    // dir exists
			"bbb/ccc/":    true,    // dir exists
			"bin/aaa.sh":  0600,    // file exists, perms match
			"bin/bbb.sh":  0600,    // file exists, perms match
			"hello.txt":   "hello", // file exists, content matches
			"unknown/":    false,   // dir should not exist
			"unknown.txt": false,   // file should not exist
		})
	})
}

func TestAssertPathsWhenNilMap(t *testing.T) {
	// Handles nil w/out error
	AssertPaths(t, "", nil)
}

func TestCurrentDir(t *testing.T) {
	cwd, _ := os.Getwd()

	assert.Equal(t, cwd, CurrentDir())
}
func TestCurrentDirWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&osGetwd, "", errors.New("boom"))
	defer stubs.Reset()

	assert.PanicsWithError(t, "boom", func() {
		CurrentDir()
	})
}

func TestInTempDir(t *testing.T) {
	startDir, _ := os.Getwd()

	InTempDir(t, func(dir string) {
		assert.DirExists(t, dir)
		cwd, _ := os.Getwd()
		dir, _ = filepath.EvalSymlinks(dir)
		assert.Equal(t, dir, cwd)
	})

	endDir, _ := os.Getwd()
	assert.Equal(t, startDir, endDir)
}

func TestInTempDirWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&osChdir, errors.New("boom"))
	defer stubs.Reset()

	// error on change into temp dir
	assert.PanicsWithError(t, "boom", func() {
		InTempDir(t, func(dir string) {})
	})

	// error on change out of temp dir
	stubs.StubFunc(&osChdir, nil)
	assert.PanicsWithError(t, "boom", func() {
		InTempDir(t, func(dir string) {
			stubs.StubFunc(&osChdir, errors.New("boom"))
		})
	})
}

func TestMkdirAll(t *testing.T) {
	InTempDir(t, func(dir string) {
		path := filepath.Join(dir, "foo", "bar", "baz")
		MkdirAll(t, path, 0755)
		assert.DirExists(t, path)
	})
}

func TestMkdirAllWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&osMkdirAll, errors.New("boom"))
	defer stubs.Reset()

	assert.PanicsWithError(t, "boom", func() {
		MkdirAll(t, "yolo", 0755)
	})
}

func TestRemoveAll(t *testing.T) {
	InTempDir(t, func(dir string) {
		path := filepath.Join(dir, "foo", "bar", "baz")

		err := os.MkdirAll(path, 0755)
		assert.NoError(t, err)
		assert.DirExists(t, path)

		RemoveAll(t, path)
		assert.NoDirExists(t, path)
	})
}

func TestRemoveAllWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&osRemoveAll, errors.New("boom"))
	defer stubs.Reset()

	assert.PanicsWithError(t, "boom", func() {
		RemoveAll(t, "yolo")
	})
}

func TestWriteFile(t *testing.T) {
	InTempDir(t, func(dir string) {
		path := filepath.Join(dir, "foo.txt")
		WriteFile(t, path, []byte(""), 0600)
		assert.FileExists(t, path)
	})
}

func TestWriteFileWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&osWriteFile, errors.New("boom"))
	defer stubs.Reset()

	assert.PanicsWithError(t, "boom", func() {
		WriteFile(t, "yolo", []byte(""), 0600)
	})
}

func TestWritePaths(t *testing.T) {
	InTempDir(t, func(dir string) {
		WritePaths(t, dir, map[string]any{
			"aaa/":       true,
			"bbb/ccc/":   true,
			"bin/aaa.sh": "aaa",
			"bin/bbb.sh": "bbb",
			"hello.txt":  "hello",
		})

		var data []byte

		assert.DirExists(t, filepath.Join(dir, "aaa"))
		assert.DirExists(t, filepath.Join(dir, "bbb", "ccc"))

		assert.FileExists(t, filepath.Join(dir, "bin", "aaa.sh"))
		data, _ = os.ReadFile(filepath.Join(dir, "bin", "aaa.sh"))
		assert.Equal(t, "aaa", string(data))

		assert.FileExists(t, filepath.Join(dir, "bin", "bbb.sh"))
		data, _ = os.ReadFile(filepath.Join(dir, "bin", "bbb.sh"))
		assert.Equal(t, "bbb", string(data))

		assert.FileExists(t, filepath.Join(dir, "hello.txt"))
		data, _ = os.ReadFile(filepath.Join(dir, "hello.txt"))
		assert.Equal(t, "hello", string(data))
	})
}

func TestWritePathsWhenNilMap(t *testing.T) {
	// Handles nil w/out error
	WritePaths(t, "", nil)
}
