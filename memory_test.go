package config_test

import (
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestMemory(t *testing.T) {
	type test struct {
		name         string
		initial      map[string]any
		transformers []config.Transformer
		expected     map[string]any
	}
	tests := []test{
		{name: "passthrough", initial: map[string]any{"hello": "world"}, transformers: nil, expected: map[string]any{"hello": "world"}},
		{name: "transform", initial: map[string]any{"hello": "world"}, transformers: []config.Transformer{
			config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				delete(m, "hello")
				m["test"] = "test"
				return m, nil
			})}, expected: map[string]any{"test": "test"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := config.NewMemory(test.initial, test.transformers...)
			ctx := &config.GetContext{}
			value, err := m.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}
			actual, ok := value.(map[string]any)
			if !ok {
				t.Fatalf("expected map but found %T", value)
			}
			for k, v := range test.expected {
				actualValue, ok := actual[k]
				if !ok {
					t.Fatalf("key '%s' not found", k)
				}
				if actualValue != v {
					t.Fatalf("expected key '%s' value '%s' to equal '%s'", k, actualValue, v)
				}
			}
			for k, v := range actual {
				_, ok := test.expected[k]
				if !ok {
					t.Fatalf("extra key '%s' value '%s' found in output", k, v)
				}
			}
		})
	}

}
