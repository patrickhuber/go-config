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
			glob: "**/config.*",
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
			glob: "**/config.*",
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
			provider := config.NewGlob(dir, test.glob)
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
