package render

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	tests := []struct {
		desc      string
		value     any
		data      map[string]any
		expected  any
		assertion assert.ErrorAssertionFunc
	}{
		// non-string scalars should get passed through unchanged.
		{desc: "[nil] noop", value: nil, expected: nil, assertion: assert.NoError},
		{desc: "[int] noop", value: 123, expected: 123, assertion: assert.NoError},
		{desc: "[float] noop", value: 1.2, expected: 1.2, assertion: assert.NoError},
		{desc: "[bool] noop", value: true, expected: true, assertion: assert.NoError},

		{
			desc:      "[string] should pass non-template strings through unchanged",
			value:     "Howdy 👋",
			expected:  "Howdy 👋",
			assertion: assert.NoError,
		},
		{
			desc:  "[string] should render template strings",
			value: "Hello, {{ .Name }} 👋",
			data: map[string]any{
				"Name": "World",
			},
			expected:  "Hello, World 👋",
			assertion: assert.NoError,
		},
		{
			desc:      "[string] should return template errors",
			value:     `{{ fail "boom" }}`,
			expected:  "",
			assertion: assert.Error,
		},

		{
			desc:      "[slice] should pass non-template values through unchanged",
			value:     []any{"111", 222, "333"},
			expected:  []any{"111", 222, "333"},
			assertion: assert.NoError,
		},
		{
			desc: "[slice] should render template strings recursively",
			value: []any{
				"{{ .aaa }}",
				[]any{
					"{{ .bbb }}",
					[]any{
						"{{ .ccc }}",
					},
				},
			},
			data: map[string]any{
				"aaa": "AAA",
				"bbb": "BBB",
				"ccc": "CCC",
			},
			expected: []any{
				"AAA",
				[]any{
					"BBB",
					[]any{
						"CCC",
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			desc:      "[slice] should return template errors",
			value:     []any{"111", "222", `{{ fail "boom" }}`},
			expected:  []any(nil),
			assertion: assert.Error,
		},

		{
			desc: "[map] should pass non-template values through unchanged",
			value: map[string]any{
				"aaa": "aaa one",
				"bbb": []any{"bbb one"},
				"ccc": map[string]any{
					"ccc.1": "ccc one",
					"ccc.2": []any{"ccc two"},
				},
			},
			expected: map[string]any{
				"aaa": "aaa one",
				"bbb": []any{"bbb one"},
				"ccc": map[string]any{
					"ccc.1": "ccc one",
					"ccc.2": []any{"ccc two"},
				},
			},
			assertion: assert.NoError,
		},
		{
			desc: "[map] should render template strings recursively",
			value: map[string]any{
				"{{ .aaa }}": "{{ .aaa }} one",
				"{{ .bbb }}": []any{"{{ .bbb }} one"},
				"{{ .ccc }}": map[string]any{
					"{{ .ccc }}.1": "{{ .ccc }} one",
					"{{ .ccc }}.2": []any{"{{ .ccc }} two"},
				},
			},
			data: map[string]any{
				"aaa": "AAA",
				"bbb": "BBB",
				"ccc": "CCC",
			},
			expected: map[string]any{
				"AAA": "AAA one",
				"BBB": []any{"BBB one"},
				"CCC": map[string]any{
					"CCC.1": "CCC one",
					"CCC.2": []any{"CCC two"},
				},
			},
			assertion: assert.NoError,
		},
		{
			desc: "[map] should return template errors from keys",
			value: map[string]any{
				`{{ fail "boom" }}`: "aaa",
			},
			expected:  map[string]any(nil),
			assertion: assert.Error,
		},
		{
			desc: "[map] should return template errors from values",
			value: map[string]any{
				"aaa": "aaa",
				"bbb": map[string]any{
					"ccc": `{{ fail "boom" }}`,
				},
			},
			expected:  map[string]any(nil),
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := Any(tt.value, tt.data)
			if tt.assertion != nil {
				tt.assertion(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// Somewhat of an integration test for the primary use case.
func TestAny_WithDataFromJSON(t *testing.T) {
	// Decode example.json into an untyped map
	var values map[string]any
	buf, err := os.ReadFile(filepath.Join("testdata", "example.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &values)
	if err != nil {
		panic(err)
	}

	// Render the map w/ some data.
	rendered, err := Any(values, map[string]any{
		"aaa": "AAA",
		"bbb": "BBB",
		"ccc": "CCC",
	})

	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"AAA": "AAA",
		"BBB": []any{
			"BBB.1",
			"BBB.2",
			[]any{
				"BBB.3.1",
				"BBB.3.2",
			},
		},
		"CCC": map[string]any{
			"CCC.1": map[string]any{
				"CCC.1.1": "CCC.1.1",
				"CCC.1.2": "CCC.1.2",
			},
		},
	}, rendered)
}

func TestRenderFile(t *testing.T) {
	tests := []struct {
		Desc     string
		Values   map[string]any
		Template string
		Rendered string
		Err      string
	}{
		{
			Desc:     "returns an error if unable to parse",
			Template: filepath.Join("testdata", "parse-error.txt"),
			Rendered: "",
			Err:      "unexpected \"}\" in command",
		},
		{
			Desc:     "returns an error if unable to render",
			Template: filepath.Join("testdata", "render-error.txt"),
			Rendered: "",
			Err:      "boom",
		},
		{
			Desc:     "does not error if template values are missing",
			Values:   map[string]any{},
			Template: filepath.Join("testdata", "valid.txt"),
			Rendered: "Hello, <no value>",
			Err:      "",
		},
		{
			Desc: "renders template values successfully",
			Values: map[string]any{
				"Name": "World",
			},
			Template: filepath.Join("testdata", "valid.txt"),
			Rendered: "Hello, World",
			Err:      "",
		},
		{
			Desc: "renders template functions successfully",
			Values: map[string]any{
				"Name": "Some Name",
			},
			Template: filepath.Join("testdata", "valid-func.txt"),
			Rendered: "Hello, some-name",
			Err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			actual, err := File(tt.Template, tt.Values)
			assert.Equal(t, tt.Rendered, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}

func TestRenderString(t *testing.T) {
	tests := []struct {
		Desc     string
		Values   map[string]any
		Template string
		Rendered string
		Err      string
	}{
		{
			Desc:     "returns an error if unable to parse",
			Template: "{{}",
			Rendered: "",
			Err:      "unexpected \"}\" in command",
		},
		{
			Desc:     "returns an error if unable to render",
			Template: "{{ fail \"boom\" }}",
			Rendered: "",
			Err:      "boom",
		},
		{
			Desc:     "does not error if template values are missing",
			Values:   map[string]any{},
			Template: "Hello, {{ .NotFound }}",
			Rendered: "Hello, <no value>",
			Err:      "",
		},
		{
			Desc: "renders template values successfully",
			Values: map[string]any{
				"Name": "World",
			},
			Template: "Hello, {{ .Name }}",
			Rendered: "Hello, World",
			Err:      "",
		},
		{
			Desc: "renders template functions successfully",
			Values: map[string]any{
				"Name": "Some Name",
			},
			Template: "Hello, {{ .Name | dasherize }}",
			Rendered: "Hello, some-name",
			Err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Desc, func(t *testing.T) {
			actual, err := String(tt.Template, tt.Values)
			assert.Equal(t, tt.Rendered, actual)
			if tt.Err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.Err)
			}
		})
	}
}
