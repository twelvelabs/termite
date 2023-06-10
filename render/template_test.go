package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCompile(t *testing.T) {
	_, err := Compile(``)
	assert.NoError(t, err)

	_, err = Compile(`{{}`)
	assert.ErrorContains(t, err, `unexpected "}" in command`)
}

func TestMustCompile(t *testing.T) {
	_ = MustCompile(``)

	assert.Panics(t, func() {
		_ = MustCompile(`{{}`)
	})
}

func TestTemplate_Render(t *testing.T) {
	ts, _ := Compile(`Hello, {{ .Name }}`)
	rendered, err := ts.Render(map[string]any{
		"Name": "World",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)

	ts, _ = Compile(`{{ .Something }}`)
	rendered, err = ts.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "", rendered)

	ts, _ = Compile(`Hello, {{ fail "boom" }}`)
	rendered, err = ts.Render(nil)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, "", rendered)
}

func TestTemplate_RenderBool(t *testing.T) {
	ts, _ := Compile(`{{ .Value }}`)
	rendered, err := ts.RenderBool(map[string]any{
		"Value": "true",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, rendered)

	ts, _ = Compile(`{{ .Value }}`)
	rendered, err = ts.RenderBool(nil)
	assert.NoError(t, err)
	assert.Equal(t, false, rendered)

	ts, _ = Compile(`{{ .Value }}`)
	rendered, err = ts.RenderBool(map[string]any{
		"Value": "not-a-bool",
	})
	assert.ErrorContains(t, err, "invalid syntax")
	assert.Equal(t, false, rendered)

	ts, _ = Compile(`{{ fail "boom" }}`)
	rendered, err = ts.RenderBool(nil)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, false, rendered)
}

func TestTemplate_RenderInt(t *testing.T) {
	ts, _ := Compile(`{{ .Value }}`)
	rendered, err := ts.RenderInt(map[string]any{
		"Value": "12",
	})
	assert.NoError(t, err)
	assert.Equal(t, 12, rendered)

	ts, _ = Compile(`{{ .Value }}`)
	rendered, err = ts.RenderInt(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, rendered)

	ts, _ = Compile(`{{ .Value }}`)
	rendered, err = ts.RenderInt(map[string]any{
		"Value": "not-a-number",
	})
	assert.ErrorContains(t, err, "unable to cast")
	assert.Equal(t, 0, rendered)

	ts, _ = Compile(`{{ fail "boom" }}`)
	rendered, err = ts.RenderInt(nil)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, 0, rendered)
}

func TestTemplate_RenderRequired(t *testing.T) {
	ts, _ := Compile(`{{ .Value }}`)
	rendered, err := ts.RenderRequired(map[string]any{
		"Value": "Howdy",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Howdy", rendered)

	ts, _ = Compile(`{{ fail "boom" }}`)
	rendered, err = ts.RenderRequired(nil)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, "", rendered)

	ts, _ = Compile(`{{ .Value }}`)
	rendered, err = ts.RenderRequired(map[string]any{
		"Value": "",
	})
	assert.ErrorContains(t, err, "empty string")
	assert.Equal(t, "", rendered)
}

func TestTemplate_MarshalText(t *testing.T) {
	ts, err := Compile(`Hello, {{ .Name }}`)
	assert.NoError(t, err)

	m, err := ts.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, `Hello, {{ .Name }}`, string(m))
}

func TestTemplate_UnmarshalText(t *testing.T) {
	ts := &Template{}
	err := ts.UnmarshalText([]byte(`Hello, {{ .Name }}`))
	assert.NoError(t, err)

	m, _ := ts.MarshalText()
	assert.Equal(t, "Hello, {{ .Name }}", string(m))

	rendered, err := ts.Render(map[string]any{
		"Name": "World",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)
}

func TestTemplate_UnmarshalFromYAML(t *testing.T) {
	s := `greeting: "Hello, {{ .Name }}"`
	mapping := struct {
		Greeting Template `yaml:"greeting"`
		Missing  Template `yaml:"missing"`
	}{}
	err := yaml.Unmarshal([]byte(s), &mapping)
	assert.NoError(t, err)

	data := map[string]any{
		"Name": "World",
	}
	rendered, err := mapping.Greeting.Render(data)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)

	rendered, err = mapping.Missing.Render(data)
	assert.NoError(t, err)
	assert.Equal(t, "", rendered)
}
