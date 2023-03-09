package render

import (
	"text/template"
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

// Render renders the template using data.
func (ts *Template) Render(data any) (string, error) {
	return execute(ts.t, data)
}
