package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestToml(t *testing.T) {
	dir := t.TempDir()
	type test struct {
		name         string
		file         string
		content      string
		transformers []config.Transformer
		expected     any
	}
	cases := []test{
		{"object", "object.toml", `str="string"`, nil, map[string]any{"str": "string"}},
		{"transform", "transform.toml", `hello="world"`, []config.Transformer{
			config.FuncTransformer(func(a any) (any, error) {
				aMap, ok := a.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("expected map[string]any but found %T", a)
				}
				delete(aMap, "hello")
				aMap["str"] = "string"
				return aMap, nil
			}),
		}, map[string]any{"str": "string"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.file)
			err := os.WriteFile(path, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}
			p := config.NewToml(path, config.FileOption{Transformers: test.transformers})
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
