package config_test

import (
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestMerge(t *testing.T) {
	type test struct {
		name     string
		from     any
		to       any
		expected any
	}
	tests := []test{
		{
			name:     "different_keys",
			from:     map[string]any{"hello": "world"},
			to:       map[string]any{"from": "here"},
			expected: map[string]any{"hello": "world", "from": "here"},
		},
		{
			name:     "same_keys",
			from:     map[string]any{"hello": "world"},
			to:       map[string]any{"hello": "here"},
			expected: map[string]any{"hello": "here"},
		},
		{
			name:     "bool",
			from:     true,
			to:       false,
			expected: false,
		},
		{
			name:     "float64",
			from:     0.1,
			to:       0.2,
			expected: 0.2,
		},
		{
			name:     "string",
			from:     "hello",
			to:       "world",
			expected: "world",
		},
		{
			name:     "slice",
			from:     []any{4.0, 5.0, 6.0},
			to:       []any{1.0, 2.0, 3.0},
			expected: []any{4.0, 5.0, 6.0, 1.0, 2.0, 3.0},
		},
		{
			name: "complex_object_slice",
			from: []any{
				map[string]any{"hello": "world"},
			},
			to: []any{
				map[string]any{"foo": "bar"},
			},
			expected: []any{
				map[string]any{"hello": "world"},
				map[string]any{"foo": "bar"},
			},
		},
		{
			name: "nested",
			from: map[string]any{
				"hello": map[string]any{"world": map[string]any{}}},
			to: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here"}},
			},
			expected: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here"}},
			},
		},
		{
			name: "nested_reverse",
			from: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here"}},
			},
			to: map[string]any{
				"hello": map[string]any{"world": map[string]any{}}},
			expected: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here"}},
			},
		},
		{
			name: "nested_existing",
			from: map[string]any{
				"hello": map[string]any{"world": map[string]any{"go": "home"}}},
			to: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here"}},
			},
			expected: map[string]any{
				"hello": map[string]any{"world": map[string]any{"from": "here", "go": "home"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := config.Merge(test.from, test.to)
			if err != nil {
				t.Fatal(err)
			}
			diff, err := config.Diff(actual, test.expected)
			if err != nil {
				t.Fatal(err)
			}
			if len(diff) != 0 {
				t.Fatalf("found differences between merged structs")
			}
		})
	}
}
