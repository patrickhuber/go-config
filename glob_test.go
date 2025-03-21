package config_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestGlob(t *testing.T) {
	type testFile struct {
		name    string
		content string
	}
	type test struct {
		name     string
		files    []testFile
		expected map[string]any
		glob     string
		resolver config.GlobProviderResolver
	}
	tests := []test{
		{
			name: "flat",
			files: []testFile{
				{
					name:    "config.yml",
					content: "yaml: test",
				},
				{
					name:    "config.json",
					content: `{"json": "test"}`,
				},
				{
					name:    "config.toml",
					content: `toml="test"`,
				},
				{
					name:    "config.env",
					content: "dotenv=test",
				},
			},
			expected: map[string]any{
				"yaml":   "test",
				"json":   "test",
				"toml":   "test",
				"dotenv": "test",
			},
			glob:     "**/config.*",
			resolver: nil,
		},
		{
			name: "level",
			files: []testFile{
				{
					name:    "./config.yml",
					content: "yaml: test",
				},
				{
					name:    "./child/config.json",
					content: `{"json": "test"}`,
				},
				{
					name:    "./child/grand/config.toml",
					content: `toml="test"`,
				},
			},
			expected: map[string]any{
				"yaml": "test",
				"json": "test",
				"toml": "test",
			},
			glob:     "**/config.*",
			resolver: nil,
		},
		{
			name: "level_resolver",
			files: []testFile{
				{
					name:    "./config.yml",
					content: "yaml: test",
				},
				{
					name:    "./child/config.json",
					content: `{"json": "test"}`,
				},
				{
					name:    "./child/grand/config.toml",
					content: `toml="test"`,
				},
			},
			expected: map[string]any{
				"yaml": "yaml",
				"json": "json",
				"toml": "toml",
			},
			glob: "**/config.*",
			resolver: func(match string) config.Provider {
				transformer := config.FuncTransformer(func(a any) (any, error) {
					aMap, ok := a.(map[string]any)
					if !ok {
						return nil, fmt.Errorf("expected map[string]any but found %T", a)
					}
					for k := range aMap {
						aMap[k] = k
					}
					return aMap, nil
				})
				ext := filepath.Ext(match)
				switch ext {
				case ".json":
					return config.NewJson(match, transformer)
				case ".yaml", ".yml":
					return config.NewYaml(match, transformer)
				case ".toml":
					return config.NewToml(match, transformer)
				}
				return nil
			},
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			dir := t.TempDir()
			for _, file := range test.files {

				fileName := filepath.Join(dir, file.name)
				fileDirectory := filepath.Dir(fileName)

				_, err := os.Stat(fileDirectory)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						err = os.MkdirAll(fileDirectory, 0666)
						if err != nil {
							t.Fatal(err)
						}
					} else {
						t.Fatal(err)
					}
				}

				err = os.WriteFile(fileName, []byte(file.content), 0666)
				if err != nil {
					t.Fatal(err)
				}
			}
			var provider config.Provider
			if test.resolver != nil {
				provider = config.NewGlobWithResolver(dir, test.glob, test.resolver)
			} else {
				provider = config.NewGlob(dir, test.glob)
			}

			ctx := &config.GetContext{}
			actual, err := provider.Get(ctx)

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
