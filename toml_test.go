package config_test

import (
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
			config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				delete(m, "hello")
				m["str"] = "string"
				return m, nil
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
