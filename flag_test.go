package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestFlag(t *testing.T) {
	type test struct {
		name      string
		args      []string
		flags     []config.Flag
		transform config.Transformer
		expected  map[string]any
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
		{
			name:  "transform",
			args:  []string{"--test", "abc"},
			flags: []config.Flag{&config.StringFlag{Name: "test", Default: "", Usage: "uses the test"}},
			transform: config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				return map[string]any{"root": m}, nil
			}),
			expected: map[string]any{"root": map[string]any{"test": "abc"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := []config.FlagOption{}
			if test.transform != nil {
				options = append(options, config.FlagOption{Transformers: []config.Transformer{test.transform}})
			}
			factory := config.NewFlag(test.flags, test.args, options...)
			ctx := &config.GetContext{}

			builder := config.NewBuilder(factory)
			root, err := builder.Build()
			if err != nil {
				t.Fatal(err)
			}

			cfg, err := root.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg, test.expected) {
				t.Fatalf("expected configurations to be equal")
			}
		})
	}
}
