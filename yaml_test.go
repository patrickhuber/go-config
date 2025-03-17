package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestYaml(t *testing.T) {
	dir := t.TempDir()
	type test struct {
		name     string
		file     string
		content  string
		expected any
	}
	cases := []test{
		{"string", "string.json", `"string"`, "string"},
		{"integer", "int.json", "1234", 1234},
		{"flat", "float.json", "1.24", float64(1.24)},
		{"boolean", "bool.json", "true", true},
		{"object", "object.json", `{"key": "value"}`, map[string]any{"key": "value"}},
		{"mobject", "mobject.json", `key: value`, map[string]any{"key": "value"}},
		{"array", "array.json", `["one", "two", "three"]`, []any{"one", "two", "three"}},
		{"marray", "marray.json", "- one\r\n- two\r\n- three", []any{"one", "two", "three"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.file)
			err := os.WriteFile(path, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}
			p := config.NewYaml(path)
			context := config.GetContext{}
			actual, err := p.Get(context)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatal("expected objects to be equal")
			}
		})
	}
}
