package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

func TestToml(t *testing.T) {
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

			f := config.NewToml(filesystem, testFilePath, config.FileOption{Transformers: test.transformers})

			builder := config.NewBuilder(f)
			root, err := builder.Build()
			if err != nil {
				t.Fatal(err)
			}

			ctx := &config.GetContext{}
			actual, err := root.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
