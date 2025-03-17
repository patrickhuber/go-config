package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestToml(t *testing.T) {
	dir := t.TempDir()
	type test struct {
		name     string
		file     string
		content  string
		expected any
	}
	cases := []test{
		{"object", "object.json", `str="string"`, map[string]any{"str": "string"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.file)
			err := os.WriteFile(path, []byte(test.content), 0666)
			if err != nil {
				t.Fatal(err)
			}
			p := config.NewToml(path)
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
