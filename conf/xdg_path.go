package conf

import (
	"path/filepath"

	"github.com/adrg/xdg" // spell: disable-line
)

var (
	// for stubbing
	xdgConfigFile = xdg.ConfigFile
)

// XDGPath returns the [XDG] config path for the app name.
//
// [XDG]: https://github.com/adrg/xdg
func XDGPath(name string) string {
	return XDGPathFor(name, "config.yaml")
}

// XDGPathFor returns the [XDG] config path for the app name and file.
//
// [XDG]: https://github.com/adrg/xdg
func XDGPathFor(app string, file string) string {
	path, err := xdgConfigFile(filepath.Join(app, file))
	if err != nil {
		panic(err)
	}
	return path
}
