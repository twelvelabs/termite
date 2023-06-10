package render

import (
	"fmt"
	"text/template"

	"github.com/spf13/cast"
)

const (
	strEmpty   string = ""
	strNoValue string = "<no value>"
)

// Compile parses a template string and returns, if successful,
// a new Template that can be rendered.
func Compile(s string) (*Template, error) {
	t, err := template.New("render.Template").Funcs(FuncMap).Parse(s)
	if err != nil {
		return &Template{}, err
	}
	return &Template{s: s, t: t}, nil
}

// MustCompile behaves like Compile, but panics on error.
func MustCompile(s string) *Template {
	t, err := Compile(s)
	if err != nil {
		panic(err)
	}
	return t
}

// Template is a lightweight wrapper for template strings.
//
// Since it implements the [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler] interfaces, it can be used with
// JSON or YAML fields containing template strings.
type Template struct {
	s string
	t *template.Template
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (ts *Template) MarshalText() ([]byte, error) {
	return []byte(ts.s), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (ts *Template) UnmarshalText(text []byte) error {
	tsp, err := Compile(string(text))
	*ts = *tsp
	return err
}

// Render renders the template as a string.
func (ts *Template) Render(data any) (string, error) {
	rendered, err := execute(ts.t, data)
	if err != nil {
		return strEmpty, err
	}
	if rendered == strEmpty || rendered == strNoValue {
		return strEmpty, nil
	}
	return rendered, nil
}

// Render renders the template as a bool.
func (ts *Template) RenderBool(data any) (bool, error) {
	rendered, err := ts.Render(data)
	if err != nil {
		return false, err
	}
	if rendered == strEmpty {
		return false, nil
	}
	return cast.ToBoolE(rendered)
}

// Render renders the template as an int.
func (ts *Template) RenderInt(data any) (int, error) {
	rendered, err := ts.Render(data)
	if err != nil {
		return 0, err
	}
	if rendered == strEmpty {
		return 0, nil
	}
	return cast.ToIntE(rendered)
}

// RenderRequired renders the template as a string,
// but returns an error if the result is empty.
func (ts *Template) RenderRequired(data any) (string, error) {
	rendered, err := ts.Render(data)
	if err != nil {
		return strEmpty, err
	}
	if rendered == strEmpty {
		return strEmpty, fmt.Errorf("evaluated to an empty string")
	}
	return rendered, nil
}
