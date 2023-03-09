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

// File renders the file at path with data.
func File(path string, data any) (string, error) {
	name := filepath.Base(path)
	t, err := template.New(name).Funcs(FuncMap).ParseFiles(path)
	if err != nil {
		return "", err
	}
	return execute(t, data)
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
	buf := bytes.Buffer{}
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
