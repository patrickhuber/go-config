package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

func TestYaml(t *testing.T) {
	// Use Target for cross-platform abstractions
	target := cross.NewTest(platform.Linux, arch.AMD64)

	// Use the filesystem from the target
	filesystem := target.FS()
	path := target.Path()

	// Use a base directory in the memory filesystem
	testDir := "/test"

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
		{"transform", "transform.yaml", `key: value`, []config.Transformer{
			config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				delete(m, "key")
				m["hello"] = "world"
				return m, nil
			})}, map[string]any{"hello": "world"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			testFilePath := path.Join(testDir, test.file)
			fileDirectory := path.Dir(testFilePath)

			// Ensure directory exists
			exists, err := filesystem.Exists(fileDirectory)
			if err != nil {
				t.Fatal(err)
			}
			if !exists {
				err := filesystem.MkdirAll(fileDirectory, 0666)
				if err != nil {
					t.Fatal(err)
				}
			}

			err = filesystem.WriteFile(testFilePath, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}

			p := config.NewYaml(filesystem, testFilePath, config.FileOption{Transformers: test.transformers})
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
