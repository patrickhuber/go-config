package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
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
			resolver: func(filesystem fs.FS, match string) config.Provider {
				transformer := config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
					for k := range m {
						m[k] = k
					}
					return m, nil
				})
				transformers := []config.Transformer{transformer}
				// Use a target just to get the file extension
				target := cross.NewTest(platform.Linux, arch.AMD64)
				ext := target.Path().Ext(match)
				switch ext {
				case ".json":
					return config.NewJson(filesystem, match, config.FileOption{Transformers: transformers})
				case ".yaml", ".yml":
					return config.NewYaml(filesystem, match, config.FileOption{Transformers: transformers})
				case ".toml":
					return config.NewToml(filesystem, match, config.FileOption{Transformers: transformers})
				}
				return nil
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use Target for cross-platform abstractions
			target := cross.NewTest(platform.Linux, arch.AMD64)

			// Create a memory-based filesystem using go-cross
			filesystem := target.FS()
			path := target.Path()

			// Use a base directory in the memory filesystem
			testDir := "/test"

			// Set up the files using the memory filesystem
			for _, file := range test.files {
				fileName := path.Join(testDir, file.name)
				fileDirectory := path.Dir(fileName)

				// Check if directory exists and create if needed using the memory filesystem
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

				err = filesystem.WriteFile(fileName, []byte(file.content), 0666)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Use the same memory filesystem for the provider
			provider := config.NewGlob(filesystem, path, testDir, test.glob, config.GlobOption{Resolver: test.resolver})

			ctx := &config.GetContext{}
			actual, err := provider.Get(ctx)

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Logf("Expected: %+v", test.expected)
				t.Logf("Actual: %+v", actual)
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
