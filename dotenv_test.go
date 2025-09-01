package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

func TestDotEnv(t *testing.T) {
	type test struct {
		name     string
		data     string
		expected map[string]any
	}
	tests := []test{
		{
			name: "kv",
			data: "KEY=VALUE",
			expected: map[string]any{
				"KEY": "VALUE",
			},
		},
		{
			name: "quoted",
			data: `KEY="VALUE"`,
			expected: map[string]any{
				"KEY": "VALUE",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use Target for cross-platform abstractions
			target := cross.NewTest(platform.Linux, arch.AMD64)

			// Use the filesystem from the target
			filesystem := target.FS()
			path := target.Path()

			// Use a base directory in the memory filesystem
			testDir := "/test"
			filePath := path.Join(testDir, ".env."+test.name)
			fileDirectory := path.Dir(filePath)

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

			err = filesystem.WriteFile(filePath, []byte(test.data), 0666)
			if err != nil {
				t.Fatal(err)
			}

			ctx := &config.GetContext{}
			providers, err := config.NewDotEnv(filesystem, filePath).Providers()
			if err != nil {
				t.Fatal(err)
			}
			if len(providers) != 1 {
				t.Fatal("expected exactly one provider")
			}
			provider := providers[0]
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
