package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestBuilder(t *testing.T) {
	type test struct {
		name                string
		additionalProviders []config.Provider
		initial             map[string]any
		expected            map[string]any
	}
	tests := []test{
		{
			name:                "passthrough",
			additionalProviders: nil,
			initial:             map[string]any{"hello": "world"},
			expected:            map[string]any{"hello": "world"},
		},
		{
			name: "transform",
			additionalProviders: []config.Provider{config.TransformProvider(func(a any) (any, error) {
				aMap, ok := a.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("expected input to be a map")
				}
				delete(aMap, "hello")
				aMap["test"] = "test"
				return aMap, nil
			})},
			initial:  map[string]any{"hello": "world"},
			expected: map[string]any{"test": "test"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			providers := []config.Provider{&FakeProvider{
				Data: test.initial,
			}}
			providers = append(providers, test.additionalProviders...)
			builder := config.NewBuilder(providers...)

			root := builder.Build()

			actual, err := root.Get(&config.GetContext{})
			if err != nil {
				t.Fatal(err)
			}
			actualMap, ok := actual.(map[string]any)
			if !ok {
				t.Fatalf("expected map but found %T", actual)
			}
			for expectedKey, expectedValue := range test.expected {
				actualValue, ok := actualMap[expectedKey]
				if !ok {
					t.Fatalf("expected to find key '%s' in result", expectedKey)
				}
				if !reflect.DeepEqual(expectedValue, actualValue) {
					t.Fatalf("expected %v to equal %v", expectedValue, actualValue)
				}
			}

			for actualKey, actualValue := range actualMap {
				_, ok := test.expected[actualKey]
				if !ok {
					t.Fatalf("found unexpected key %s and value %s in result", actualKey, actualValue)
				}
			}
		})
	}
}

type FakeProvider struct {
	Data  any
	Error error
}

func (p *FakeProvider) Get(ctx *config.GetContext) (any, error) {
	if p.Error != nil {
		return nil, p.Error
	}
	return p.Data, nil
}
