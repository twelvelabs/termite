package conf

import (
	"errors"
	"os"
	"reflect"

	"github.com/caarlos0/env/v8"  // spell: disable-line
	"github.com/creasty/defaults" // spell: disable-line
	yaml "gopkg.in/yaml.v3"

	"github.com/twelvelabs/termite/validate"
)

var (
	// for stubbing
	defaultsSet   = defaults.Set
	osReadFile    = os.ReadFile
	yamlUnmarshal = yaml.Unmarshal
)

// NewLoader returns a new [Loader].
func NewLoader[C any](config C, path string) *Loader[C] {
	v := reflect.ValueOf(config)
	if v.Type().Kind() != reflect.Ptr {
		panic("config struct must be a pointer")
	}
	return &Loader[C]{
		Config: config,
		Path:   path,
	}
}

// Loader populates and validates config data from a file.
type Loader[C any] struct {
	Config C
	Path   string

	loaded bool
}

// Load reads the config file at Path into the Config struct and returns it.
// Subsequent calls will bypass the read and just return the cached struct.
func (l *Loader[C]) Load() (C, error) {
	if l.loaded {
		return l.Config, nil
	}
	return l.load()
}

// Reload forces a reread of the config file at Path.
func (l *Loader[C]) Reload() (C, error) {
	return l.load()
}

func (l *Loader[C]) load() (C, error) {
	l.loaded = true
	// Set defaults
	if err := defaultsSet(l.Config); err != nil {
		return l.Config, err
	}
	// Try to read the file
	bytes, err := osReadFile(l.Path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Probably a permissions error
		return l.Config, err
	} else if err == nil {
		// Unmarshal
		err = yamlUnmarshal(bytes, l.Config)
		if err != nil {
			return l.Config, err
		}
	}
	// Override values passed in via ENV var
	if err := env.Parse(l.Config); err != nil {
		return l.Config, err
	}
	// Validate
	if err := validate.Struct(l.Config); err != nil {
		return l.Config, err
	}
	return l.Config, nil
}
