package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestFlag(t *testing.T) {
	type test struct {
		name     string
		args     []string
		flags    []config.Flag
		expected map[string]any
	}

	tests := []test{
		{
			name:     "string",
			args:     []string{"--test", "abc"},
			flags:    []config.Flag{&config.StringFlag{Name: "test", Default: "", Usage: "uses the test"}},
			expected: map[string]any{"test": "abc"},
		},
		{
			name:     "repeat",
			args:     []string{"--test", "abc", "--test", "123"},
			flags:    []config.Flag{&config.StringSliceFlag{Name: "test", Default: nil, Usage: "uses the test"}},
			expected: map[string]any{"test": []any{"abc", "123"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := config.NewFlag(test.flags, test.args)
			ctx := &config.GetContext{}
			cfg, err := p.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg, test.expected) {
				t.Fatalf("expected configurations to be equal")
			}
		})
	}
}
