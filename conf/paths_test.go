package conf

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func TestConfigDir(t *testing.T) {
	tests := []struct {
		desc            string
		HOME            string
		AppData         string
		XDG_CONFIG_HOME string
		windows         bool
		expected        string
	}{
		{
			desc:            "[non-windows] default",
			HOME:            "HOME_DIR",
			AppData:         "APP_DATA_DIR",
			XDG_CONFIG_HOME: "",
			windows:         false,
			expected:        filepath.Join("HOME_DIR", ".config", "my-app"),
		},
		{
			desc:            "[non-windows] xdg",
			HOME:            "HOME_DIR",
			AppData:         "APP_DATA_DIR",
			XDG_CONFIG_HOME: "XDG_DIR",
			windows:         false,
			expected:        filepath.Join("XDG_DIR", "my-app"),
		},

		{
			desc:            "[windows] default",
			HOME:            "HOME_DIR",
			AppData:         "",
			XDG_CONFIG_HOME: "",
			windows:         true,
			expected:        filepath.Join("HOME_DIR", ".config", "my-app"),
		},
		{
			desc:            "[windows] app data",
			HOME:            "HOME_DIR",
			AppData:         "APP_DATA_DIR",
			XDG_CONFIG_HOME: "",
			windows:         true,
			expected:        filepath.Join("APP_DATA_DIR", "my-app"),
		},
		{
			desc:            "[windows] xdg",
			HOME:            "HOME_DIR",
			AppData:         "APP_DATA_DIR",
			XDG_CONFIG_HOME: "XDG_DIR",
			windows:         true,
			expected:        filepath.Join("XDG_DIR", "my-app"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			stubs := gostub.StubFunc(&isWindowsFunc, tt.windows)
			defer stubs.Reset()

			t.Setenv("HOME", tt.HOME)
			t.Setenv("AppData", tt.AppData)
			t.Setenv("XDG_CONFIG_HOME", tt.XDG_CONFIG_HOME)

			assert.Equal(t, tt.expected, ConfigDir("my-app"))
		})
	}
}

func TestStateDir_WhenWindows(t *testing.T) {
	tests := []struct {
		desc           string
		HOME           string
		LocalAppData   string
		XDG_STATE_HOME string
		windows        bool
		expected       string
	}{
		{
			desc:           "[non-windows] default",
			HOME:           "HOME_DIR",
			LocalAppData:   "APP_DATA_DIR",
			XDG_STATE_HOME: "",
			windows:        false,
			expected:       filepath.Join("HOME_DIR", ".local", "state", "my-app"),
		},
		{
			desc:           "[non-windows] xdg",
			HOME:           "HOME_DIR",
			LocalAppData:   "APP_DATA_DIR",
			XDG_STATE_HOME: "XDG_DIR",
			windows:        false,
			expected:       filepath.Join("XDG_DIR", "my-app"),
		},

		{
			desc:           "[windows] default",
			HOME:           "HOME_DIR",
			LocalAppData:   "",
			XDG_STATE_HOME: "",
			windows:        true,
			expected:       filepath.Join("HOME_DIR", ".local", "state", "my-app"),
		},
		{
			desc:           "[windows] app data",
			HOME:           "HOME_DIR",
			LocalAppData:   "APP_DATA_DIR",
			XDG_STATE_HOME: "",
			windows:        true,
			expected:       filepath.Join("APP_DATA_DIR", "my-app"),
		},
		{
			desc:           "[windows] xdg",
			HOME:           "HOME_DIR",
			LocalAppData:   "APP_DATA_DIR",
			XDG_STATE_HOME: "XDG_DIR",
			windows:        true,
			expected:       filepath.Join("XDG_DIR", "my-app"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			stubs := gostub.StubFunc(&isWindowsFunc, tt.windows)
			defer stubs.Reset()

			t.Setenv("HOME", tt.HOME)
			t.Setenv("LocalAppData", tt.LocalAppData)
			t.Setenv("XDG_STATE_HOME", tt.XDG_STATE_HOME)

			assert.Equal(t, tt.expected, StateDir("my-app"))
		})
	}
}

func TestDataDir_WhenWindows(t *testing.T) {
	tests := []struct {
		desc          string
		HOME          string
		LocalAppData  string
		XDG_DATA_HOME string
		windows       bool
		expected      string
	}{
		{
			desc:          "[non-windows] default",
			HOME:          "HOME_DIR",
			LocalAppData:  "APP_DATA_DIR",
			XDG_DATA_HOME: "",
			windows:       false,
			expected:      filepath.Join("HOME_DIR", ".local", "share", "my-app"),
		},
		{
			desc:          "[non-windows] xdg",
			HOME:          "HOME_DIR",
			LocalAppData:  "APP_DATA_DIR",
			XDG_DATA_HOME: "XDG_DIR",
			windows:       false,
			expected:      filepath.Join("XDG_DIR", "my-app"),
		},

		{
			desc:          "[windows] default",
			HOME:          "HOME_DIR",
			LocalAppData:  "",
			XDG_DATA_HOME: "",
			windows:       true,
			expected:      filepath.Join("HOME_DIR", ".local", "share", "my-app"),
		},
		{
			desc:          "[windows] app data",
			HOME:          "HOME_DIR",
			LocalAppData:  "APP_DATA_DIR",
			XDG_DATA_HOME: "",
			windows:       true,
			expected:      filepath.Join("APP_DATA_DIR", "my-app"),
		},
		{
			desc:          "[windows] xdg",
			HOME:          "HOME_DIR",
			LocalAppData:  "APP_DATA_DIR",
			XDG_DATA_HOME: "XDG_DIR",
			windows:       true,
			expected:      filepath.Join("XDG_DIR", "my-app"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			stubs := gostub.StubFunc(&isWindowsFunc, tt.windows)
			defer stubs.Reset()

			t.Setenv("HOME", tt.HOME)
			t.Setenv("LocalAppData", tt.LocalAppData)
			t.Setenv("XDG_DATA_HOME", tt.XDG_DATA_HOME)

			assert.Equal(t, tt.expected, DataDir("my-app"))
		})
	}
}

func TestIsWindows(t *testing.T) {
	// Got to get that coverage, yo.
	assert.Equal(t, (runtime.GOOS == "windows"), isWindows())
}
