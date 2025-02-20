package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestDiff(t *testing.T) {
	type test struct {
		name     string
		from     any
		to       any
		expected []config.Change
	}
	tests := []test{
		{
			name: "string_update",
			from: "hello",
			to:   "world",
			expected: []config.Change{
				{
					ChangeType: config.Update,
					From:       "hello",
					To:         "world",
				},
			},
		},
		{
			name: "float_update",
			from: float64(1),
			to:   float64(2),
			expected: []config.Change{
				{
					ChangeType: config.Update,
					From:       float64(1),
					To:         float64(2),
				},
			},
		},
		{
			name: "bool_update",
			from: true,
			to:   false,
			expected: []config.Change{
				{
					ChangeType: config.Update,
					From:       true,
					To:         false,
				},
			},
		},
		{
			name: "map_delete",
			from: map[string]any{"hello": "world"},
			to:   map[string]any{},
			expected: []config.Change{
				{
					Path:       []string{"hello"},
					ChangeType: config.Delete,
					From:       "world",
					To:         nil,
				},
			},
		},
		{
			name: "map_create",
			from: map[string]any{},
			to:   map[string]any{"hello": "world"},
			expected: []config.Change{
				{
					Path:       []string{"hello"},
					ChangeType: config.Create,
					From:       nil,
					To:         "world",
				},
			},
		},
		{
			name: "map_update",
			from: map[string]any{"hello": "world"},
			to:   map[string]any{"hello": "config"},
			expected: []config.Change{
				{
					Path:       []string{"hello"},
					ChangeType: config.Update,
					From:       "world",
					To:         "config",
				},
			},
		},
		{
			name: "map_and_not_map",
			from: map[string]any{"hello": "world"},
			to:   "string",
			expected: []config.Change{
				{
					Path:       nil,
					ChangeType: config.Update,
					From:       map[string]any{"hello": "world"},
					To:         "string",
				},
			},
		},
		{
			name: "not_map_and_map",
			from: "string",
			to:   map[string]any{"hello": "world"},
			expected: []config.Change{
				{
					Path:       nil,
					ChangeType: config.Update,
					From:       "string",
					To:         map[string]any{"hello": "world"},
				},
			},
		},
		{
			name: "slice_delete_0",
			from: []any{"hello", "world"},
			to:   []any{"world"},
			expected: []config.Change{
				{
					Path:       []string{"0"},
					ChangeType: config.Update,
					From:       "hello",
					To:         "world",
				},
				{
					Path:       []string{"1"},
					ChangeType: config.Delete,
					From:       "world",
					To:         nil,
				},
			},
		},
		{
			name: "slice_create_start",
			from: []any{"world"},
			to:   []any{"hello", "world"},
			expected: []config.Change{
				{
					Path:       []string{"0"},
					ChangeType: config.Update,
					From:       "world",
					To:         "hello",
				},
				{
					Path:       []string{"1"},
					ChangeType: config.Create,
					From:       nil,
					To:         "world",
				},
			},
		},
		{
			name: "slice_create_end",
			from: []any{"hello"},
			to:   []any{"hello", "world"},
			expected: []config.Change{
				{
					Path:       []string{"1"},
					ChangeType: config.Create,
					From:       nil,
					To:         "world",
				},
			},
		},
		{
			name: "slice_create_insert",
			from: []any{"hello", "world"},
			to:   []any{"hello", "new", "world"},
			expected: []config.Change{
				{
					Path:       []string{"1"},
					ChangeType: config.Update,
					From:       "world",
					To:         "new",
				},
				{
					Path:       []string{"2"},
					ChangeType: config.Create,
					From:       nil,
					To:         "world",
				},
			},
		},
		{
			name: "slice_delete_middle",
			from: []any{"hello", "new", "world"},
			to:   []any{"hello", "world"},
			expected: []config.Change{
				{
					Path:       []string{"1"},
					ChangeType: config.Update,
					From:       "new",
					To:         "world",
				},
				{
					Path:       []string{"2"},
					ChangeType: config.Delete,
					From:       "world",
					To:         nil,
				},
			},
		},
		{
			name: "new_slice",
			from: []any{"hello", "world"},
			to:   []any{"new", "elements"},
			expected: []config.Change{
				{
					Path:       []string{"0"},
					ChangeType: config.Update,
					From:       "hello",
					To:         "new",
				},
				{
					Path:       []string{"1"},
					ChangeType: config.Update,
					From:       "world",
					To:         "elements",
				},
			},
		},
		{
			name: "slice_reverse",
			from: []any{"hello", "new", "world"},
			to:   []any{"world", "new", "hello"},
			expected: []config.Change{
				{
					Path:       []string{"0"},
					ChangeType: config.Update,
					From:       "hello",
					To:         "world",
				},
				{
					Path:       []string{"2"},
					ChangeType: config.Update,
					From:       "world",
					To:         "hello",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := config.Diff(test.from, test.to)
			if err != nil {
				t.Fatal(err)
			}
			err = assertChangesEqual(test.expected, actual)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertChangesEqual(expected, actual []config.Change) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected change count of %d actual change count %d", len(expected), len(actual))
	}
	for i := range expected {
		err := assertChangeEqual(expected[i], actual[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func assertChangeEqual(expected, actual config.Change) error {
	if expected.ChangeType != actual.ChangeType {
		return fmt.Errorf("expected change type %s but found change type %s", expected.ChangeType, actual.ChangeType)
	}
	if len(expected.Path) != len(actual.Path) {
		return fmt.Errorf("expected path length to be %d but found %d", len(expected.Path), len(actual.Path))
	}
	if !reflect.DeepEqual(expected.From, actual.From) {
		return fmt.Errorf("expected To to be %v but found %v", expected.From, actual.From)
	}
	if !reflect.DeepEqual(expected.To, actual.To) {
		return fmt.Errorf("expected From to be %v but found %v", expected.To, actual.To)
	}
	return nil
}
