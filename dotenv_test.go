package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
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
			dir := t.TempDir()
			filePath := filepath.Join(dir, ".env."+test.name)

			file, err := os.Create(filePath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			_, err = file.WriteString(test.data)
			if err != nil {
				t.Fatal(err)
			}
			file.Close()

			actual, err := config.NewDotEnv(filePath).Get()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
