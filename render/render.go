package render

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobuffalo/flect"
)

var (
	// The function map to use when rendering templates.
	FuncMap = DefaultFuncMap()
)

// DefaultFuncMap returns a map of functions for the template engine.
// It includes everything from [sprig.FuncMap], plus a number
// of functions from [flect].
func DefaultFuncMap() template.FuncMap {
	funcs := sprig.FuncMap()

	// See: https://pkg.go.dev/github.com/gobuffalo/flect
	funcs["camelize"] = flect.Camelize
	funcs["capitalize"] = flect.Capitalize
	funcs["dasherize"] = flect.Dasherize
	funcs["humanize"] = flect.Humanize
	funcs["ordinalize"] = flect.Ordinalize
	funcs["pascalize"] = flect.Pascalize
	funcs["pluralize"] = flect.Pluralize
	funcs["singularize"] = flect.Singularize
	funcs["titleize"] = flect.Titleize
	funcs["underscore"] = flect.Underscore

	return funcs
}

// Any attempts to render strings in the given value with data.
// The input value is assumed to come from an untyped map[string]any
// (typically from decoding unknown JSON or YAML).
//
// Delegates internally to [Map], [Slice], and [String]
// (see the documentation for those functions for more info).
// Other types are returned unchanged.
func Any(value any, data any) (any, error) {
	switch casted := value.(type) {
	case map[string]any:
		return Map(casted, data)
	case []any:
		return Slice(casted, data)
	case string:
		return String(casted, data)
	default:
		return casted, nil
	}
}

// File renders the file at path with data.
func File(path string, data any) (string, error) {
	name := filepath.Base(path)
	t, err := template.New(name).Funcs(FuncMap).ParseFiles(path)
	if err != nil {
		return "", err
	}
	return execute(t, data)
}

// Map recursively renders the keys and values of the given map with data.
func Map(values map[string]any, data any) (map[string]any, error) {
	rendered := map[string]any{}
	for k, v := range values {
		ka, err := Any(k, data)
		if err != nil {
			return nil, err
		}
		// We know if `k` was originally a string, `Any` should return as one.
		k = ka.(string)

		va, err := Any(v, data)
		if err != nil {
			return nil, err
		}
		rendered[k] = va
	}
	return rendered, nil
}

// Map recursively renders the elements of the given slice with data.
func Slice(values []any, data any) ([]any, error) {
	var err error
	for idx := range values {
		values[idx], err = Any(values[idx], data)
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

// String renders the template string with data.
func String(s string, data any) (string, error) {
	t, err := template.New("render.String").Funcs(FuncMap).Parse(s)
	if err != nil {
		return "", err
	}
	return execute(t, data)
}

func execute(t *template.Template, data any) (string, error) {
	if t == nil {
		return strEmpty, nil
	}
	buf := bytes.Buffer{}
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
