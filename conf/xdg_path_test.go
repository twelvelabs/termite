package conf

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func ExampleXDGPath() {
	// Returns:
	//  - ~/.config/termite/config.yaml on Linux.
	//  - ~/Library/Application Support/termite/config.yaml on macOS.
	//  - %LocalAppData%\termite\config.yaml on Windows.
	XDGPath("termite")
}

func TestXDGPath(t *testing.T) {
	path, _ := xdg.ConfigFile(filepath.Join("termite", "config.yaml"))
	assert.Equal(t, path, XDGPath("termite"))
}

func ExampleXDGPathFor() {
	// Returns:
	//  - ~/.config/termite/other.yaml on Linux.
	//  - ~/Library/Application Support/termite/other.yaml on macOS.
	//  - %LocalAppData%\termite\other.yaml on Windows.
	XDGPathFor("termite", "other.yaml")
}

func TestXDGPathFor(t *testing.T) {
	path, _ := xdg.ConfigFile(filepath.Join("termite", "special.yaml"))
	assert.Equal(t, path, XDGPathFor("termite", "special.yaml"))
}

func TestXDGPathForWhenError(t *testing.T) {
	stubs := gostub.StubFunc(&xdgConfigFile, "", errors.New("boom"))
	defer stubs.Reset()

	assert.PanicsWithError(t, "boom", func() {
		XDGPathFor("termite", "nope.yaml")
	})
}
