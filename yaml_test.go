package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestYaml(t *testing.T) {
	dir := t.TempDir()
	type test struct {
		name         string
		file         string
		content      string
		transformers []config.Transformer
		expected     any
	}
	cases := []test{
		{"string", "string.yaml", `"string"`, nil, "string"},
		{"integer", "int.yaml", "1234", nil, 1234},
		{"flat", "float.yaml", "1.24", nil, float64(1.24)},
		{"boolean", "bool.yaml", "true", nil, true},
		{"object", "object.yaml", `{"key": "value"}`, nil, map[string]any{"key": "value"}},
		{"mobject", "mobject.yaml", `key: value`, nil, map[string]any{"key": "value"}},
		{"array", "array.yaml", `["one", "two", "three"]`, nil, []any{"one", "two", "three"}},
		{"marray", "marray.yaml", "- one\r\n- two\r\n- three", nil, []any{"one", "two", "three"}},
		{"transform", "transform.yaml", `key: value`, []config.Transformer{config.FuncTransformer(func(a any) (any, error) {
			aMap, ok := a.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map but found %T", a)
			}
			delete(aMap, "key")
			aMap["hello"] = "world"
			return aMap, nil
		})}, map[string]any{"hello": "world"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.file)
			err := os.WriteFile(path, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}
			p := config.NewYaml(path, config.FileOption{Transformers: test.transformers})
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
