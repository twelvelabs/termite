package conf

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	appData       = "AppData"
	localAppData  = "LocalAppData"
	userHome      = "HOME"
	xdgConfigHome = "XDG_CONFIG_HOME"
	xdgDataHome   = "XDG_DATA_HOME"
	xdgStateHome  = "XDG_STATE_HOME"
)

var isWindowsFunc = isWindows

// ConfigFile returns the default config path for the app.
func ConfigFile(app string) string {
	return ConfigFileFor(app, "config.yaml")
}

// ConfigFileFor returns the config path for the app and file name.
func ConfigFileFor(app string, filename string) string {
	return filepath.Join(ConfigDir(app), filename)
}

// ConfigDir returns the path to the config dir for the app.
// Path precedence:
//
//   - $XDG_CONFIG_HOME/$name
//   - $AppData/$name (windows only)
//   - $HOME/.config/$name
func ConfigDir(app string) string {
	var path string
	if a := os.Getenv(xdgConfigHome); a != "" {
		path = filepath.Join(a, app)
	} else if b := os.Getenv(appData); isWindowsFunc() && b != "" {
		path = filepath.Join(b, app)
	} else {
		c, _ := os.UserHomeDir()
		path = filepath.Join(c, ".config", app)
	}
	return path
}

// DataDir returns the path to the local state dir for the app.
// Path precedence:
//
//   - $XDG_DATA_HOME/$name
//   - $LocalAppData/$name (windows only)
//   - $HOME/.local/share/$name
func DataDir(app string) string {
	var path string
	if a := os.Getenv(xdgDataHome); a != "" {
		path = filepath.Join(a, app)
	} else if b := os.Getenv(localAppData); isWindowsFunc() && b != "" {
		path = filepath.Join(b, app)
	} else {
		c, _ := os.UserHomeDir()
		path = filepath.Join(c, ".local", "share", app)
	}
	return path
}

// StateDir returns the path to the local state dir for the app.
// Path precedence:
//
//   - $XDG_STATE_HOME/$name
//   - $LocalAppData/$name (windows only)
//   - $HOME/.local/state/$name
func StateDir(app string) string {
	var path string
	if a := os.Getenv(xdgStateHome); a != "" {
		path = filepath.Join(a, app)
	} else if b := os.Getenv(localAppData); isWindowsFunc() && b != "" {
		path = filepath.Join(b, app)
	} else {
		c, _ := os.UserHomeDir()
		path = filepath.Join(c, ".local", "state", app)
	}
	return path
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
