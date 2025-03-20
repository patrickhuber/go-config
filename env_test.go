package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
	"maps"
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
			[]config.Transformer{config.FuncTransformer(func(instance any) (any, error) {
				envMap := map[string]any{}
				instanceMap, ok := instance.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("expected instance to be of type map[string]any")
				}
				maps.Copy(envMap, instanceMap)
				return map[string]any{"env": envMap}, nil
			})},
			map[string]any{"env": map[string]any{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range test.env {
				err := os.Setenv(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}
			p := config.NewEnv(test.prefix, test.transforms...)
			actual, err := p.Get()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("expected configurations to match")
			}
		})
	}
}
