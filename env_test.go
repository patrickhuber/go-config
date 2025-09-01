package config_test

import (
	"reflect"
	"testing"

	"maps"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
)

func TestEnv(t *testing.T) {
	type test struct {
		name       string
		env        map[string]string
		prefix     string
		transforms []config.Transformer
		expected   map[string]any
	}
	tests := []test{
		{
			"prefix",
			map[string]string{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
			"TEST",
			nil,
			map[string]any{"TEST1": "TEST1", "TEST2": "TEST2"},
		},
		{
			"noprefix",
			map[string]string{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
			"",
			nil,
			map[string]any{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
		},
		{
			"prefix_transform",
			map[string]string{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
			"",
			[]config.Transformer{config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				envMap := map[string]any{}
				maps.Copy(envMap, m)
				return map[string]any{"env": envMap}, nil
			})},
			map[string]any{"env": map[string]any{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a memory environment with the test data
			testEnv := env.NewMemoryWithMap(test.env)

			factory := config.NewEnv(testEnv, config.EnvOption{
				Transformers: test.transforms,
				Prefix:       test.prefix})

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
				t.Fatalf("expected configurations to match")
			}
		})
	}
}
