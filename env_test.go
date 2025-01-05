package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestEnv(t *testing.T) {
	type test struct {
		name     string
		env      map[string]string
		prefix   string
		expected map[string]any
	}
	tests := []test{
		{
			"prefix",
			map[string]string{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
			"TEST",
			map[string]any{"TEST1": "TEST1", "TEST2": "TEST2"},
		},
		{
			"noprefix",
			map[string]string{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
			"",
			map[string]any{"TEST1": "TEST1", "TEST2": "TEST2", "NOTEST": "NOTEST"},
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
			p := config.NewEnv(test.prefix)
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
