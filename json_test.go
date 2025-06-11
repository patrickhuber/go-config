package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestJson(t *testing.T) {
	dir := t.TempDir()
	type test struct {
		name         string
		file         string
		content      string
		transformers []config.Transformer
		expected     any
	}
	cases := []test{
		{"string", "string.json", `"string"`, nil, "string"},
		{"integer", "int.json", "1234", nil, float64(1234)},
		{"flat", "float.json", "1.24", nil, float64(1.24)},
		{"boolean", "bool.json", "true", nil, true},
		{"object", "object.json", `{"key": "value"}`, nil, map[string]any{"key": "value"}},
		{"array", "array.json", `["one", "two", "three"]`, nil, []any{"one", "two", "three"}},
		{"transform", "transform.json", `{"key": "value"}`, []config.Transformer{
			config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				for k := range m {
					m[k] = k
				}
				return m, nil
			}),
		}, map[string]any{"key": "key"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.file)
			err := os.WriteFile(path, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}
			p := config.NewJson(path, config.FileOption{Transformers: test.transformers})
			ctx := &config.GetContext{}
			actual, err := p.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
