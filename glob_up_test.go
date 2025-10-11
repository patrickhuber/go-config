package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

func TestGlobUp(t *testing.T) {
	type testFile struct {
		name    string
		content string
	}
	type test struct {
		dir      string
		name     string
		files    []testFile
		expected map[string]any
		glob     string
	}
	tests := []test{
		{
			dir:  ".",
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
			},
			expected: map[string]any{
				"yaml": "test",
				"json": "test",
				"toml": "test",
			},
			glob: "**/config.*",
		},
		{
			dir:  "./child/grand",
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
			glob: "**/config.*",
		},
	}
	platforms := []platform.Platform{
		platform.Linux,
		platform.Windows,
		platform.Darwin,
	}

	for _, test := range tests {
		for _, plat := range platforms {
			t.Run(test.name+"_"+plat.String(), func(t *testing.T) {
				// Use Target for cross-platform abstractions
				target := cross.NewTest(plat, arch.AMD64)

				// Create a memory-based filesystem using go-cross
				filesystem := target.FS()
				path := target.Path()

				// Use a base directory in the memory filesystem
				testDir := "/test"
				if platform.IsWindows(plat) {
					testDir = `C:\test`
				}

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

				workingDirectory := path.Join(testDir, test.dir)
				resolver := config.DefaultGlobResolver(filesystem, path)
				factory := config.NewGlobUp(filesystem, path, resolver, workingDirectory, test.glob)
				builder := config.NewBuilder(factory)
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
					t.Logf("Expected: %+v", test.expected)
					t.Logf("Actual: %+v", actual)
					t.Fatal("expected objects to be equal")
				}
			})
		}
	}
}
