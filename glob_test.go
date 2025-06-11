package config_test

import (
	"errors"
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
				transformer := config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
					for k := range m {
						m[k] = k
					}
					return m, nil
				})
				transformers := []config.Transformer{transformer}
				ext := filepath.Ext(match)
				switch ext {
				case ".json":
					return config.NewJson(match, config.FileOption{Transformers: transformers})
				case ".yaml", ".yml":
					return config.NewYaml(match, config.FileOption{Transformers: transformers})
				case ".toml":
					return config.NewToml(match, config.FileOption{Transformers: transformers})
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

			provider := config.NewGlob(dir, test.glob, config.GlobOption{Resolver: test.resolver})

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
