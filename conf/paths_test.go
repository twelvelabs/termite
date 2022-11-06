package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func TestConfigDir_WhenWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, true)
	defer stubs.Reset()

	// default is ~/.config/$name
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".config", "my-app"), ConfigDir("my-app"))

	// but prefers $AppData/$name
	t.Setenv(appData, "DATA_DIR")

	fmt.Println("")
	fmt.Println(" ==> os.Environ")
	for _, line := range os.Environ() {
		fmt.Println(line)
	}
	fmt.Println("")
	fmt.Println("")

	assert.Equal(t, true, isWindowsFunc())
	assert.Equal(t, "DATA_DIR", os.Getenv(appData))

	assert.Equal(t, filepath.Join("DATA_DIR", "my-app"), ConfigDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgConfigHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), ConfigDir("my-app"))
}

func TestConfigDir_WhenNotWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, false)
	defer stubs.Reset()

	t.Setenv(appData, "DATA_DIR")

	// default is ~/.config/$name (regardless of $AppData)
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".config", "my-app"), ConfigDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgConfigHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), ConfigDir("my-app"))
}

func TestStateDir_WhenWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, true)
	defer stubs.Reset()

	// default is ~/.local/state/$name
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".local", "state", "my-app"), StateDir("my-app"))

	// but prefers $LocalAppData/$name
	t.Setenv(localAppData, "DATA_DIR")
	assert.Equal(t, filepath.Join("DATA_DIR", "my-app"), StateDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgStateHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), StateDir("my-app"))
}

func TestStateDir_WhenNotWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, false)
	defer stubs.Reset()

	t.Setenv(localAppData, "DATA_DIR")

	// default is ~/.local/state/$name
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".local", "state", "my-app"), StateDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgStateHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), StateDir("my-app"))
}

func TestDataDir_WhenWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, true)
	defer stubs.Reset()

	// default is ~/.local/state/$name
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".local", "share", "my-app"), DataDir("my-app"))

	// but prefers $LocalAppData/$name
	t.Setenv(localAppData, "DATA_DIR")
	assert.Equal(t, filepath.Join("DATA_DIR", "my-app"), DataDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgDataHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), DataDir("my-app"))
}

func TestDataDir_WhenNotWindows(t *testing.T) {
	stubs := gostub.StubFunc(&isWindowsFunc, false)
	defer stubs.Reset()

	t.Setenv(localAppData, "DATA_DIR")

	// default is ~/.local/state/$name
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".local", "share", "my-app"), DataDir("my-app"))

	// XDG dir takes precedence over the above.
	t.Setenv(xdgDataHome, "XDG_DIR")
	assert.Equal(t, filepath.Join("XDG_DIR", "my-app"), DataDir("my-app"))
}

func TestIsWindows(t *testing.T) {
	// Got to get that coverage, yo.
	assert.Equal(t, (runtime.GOOS == "windows"), isWindows())
}
