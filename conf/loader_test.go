package conf

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/prashantv/gostub" // spell: disable-line
	"github.com/stretchr/testify/assert"
)

func ExampleLoader() {
	type Config struct {
		DatabaseURL string `default:"postgres://0.0.0.0:5432/db" validate:"required,url"`
		Debug       bool   `default:"true"`
	}

	loader := NewLoader(&Config{}, XDGPath("termite"))
	config, err := loader.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DatabaseURL: %v, Debug: %v", config.DatabaseURL, config.Debug)
	// Output: DatabaseURL: postgres://0.0.0.0:5432/db, Debug: true
}

type MyConfig struct {
	Name    string `default:"untitled"`
	Count   int    `default:"1" validate:"lt=100"`
	Enabled bool
	Tags    []string
}

func FixturePath(names ...string) string {
	name := "config.yaml"
	if len(names) > 0 {
		name = names[0]
	}
	return filepath.Join("testdata", name)
}

func FixtureWithDefaults() *MyConfig {
	return &MyConfig{
		Name:    "untitled",
		Count:   1,
		Enabled: false,
		Tags:    nil,
	}
}

func Fixture() *MyConfig {
	return &MyConfig{
		Name:    "aaa",
		Count:   10,
		Enabled: true,
		Tags: []string{
			"foo",
			"bar",
			"baz",
		},
	}
}

func TestNewLoader(t *testing.T) {
	assert.PanicsWithValue(t, "config struct must be a pointer", func() {
		NewLoader(MyConfig{}, "")
	})
	_ = NewLoader(&MyConfig{}, "")
}

func TestLoader_Load(t *testing.T) {
	// Config file doesn't exist - should use defaults
	loader := NewLoader(&MyConfig{}, "nope.yaml")
	config, err := loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, FixtureWithDefaults(), config)

	// Config file exists.
	loader = NewLoader(&MyConfig{}, FixturePath())
	config, err = loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, Fixture(), config)

	stubs := gostub.StubFunc(&osReadFile, []byte{}, errors.New("boom"))
	defer stubs.Reset()

	// Config loads are cached.
	config, err = loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, Fixture(), config)
}

func TestLoader_LoadWhenDefaultsError(t *testing.T) {
	stubs := gostub.StubFunc(&defaultsSet, errors.New("boom"))
	defer stubs.Reset()

	_, err := NewLoader(&MyConfig{}, FixturePath()).Load()
	assert.ErrorContains(t, err, "boom")
}

func TestLoader_LoadWhenReadError(t *testing.T) {
	stubs := gostub.StubFunc(&osReadFile, []byte{}, errors.New("boom"))
	defer stubs.Reset()

	_, err := NewLoader(&MyConfig{}, FixturePath()).Load()
	assert.ErrorContains(t, err, "boom")
}

func TestLoader_LoadWhenYAMLError(t *testing.T) {
	stubs := gostub.StubFunc(&yamlUnmarshal, errors.New("boom"))
	defer stubs.Reset()

	_, err := NewLoader(&MyConfig{}, FixturePath()).Load()
	assert.ErrorContains(t, err, "boom")
}

func TestLoader_LoadWhenValidationError(t *testing.T) {
	_, err := NewLoader(&MyConfig{}, FixturePath("invalid.yaml")).Load()
	assert.ErrorContains(t, err, "Count must be less than 100")
}

func TestLoader_Reload(t *testing.T) {
	loader := NewLoader(&MyConfig{}, FixturePath())
	config, err := loader.Load()
	assert.NoError(t, err)
	assert.Equal(t, Fixture(), config)

	stubs := gostub.StubFunc(&osReadFile, []byte{}, errors.New("boom"))
	defer stubs.Reset()

	// Reload should attempt to reread the config path.
	_, err = loader.Reload()
	assert.ErrorContains(t, err, "boom")
}
