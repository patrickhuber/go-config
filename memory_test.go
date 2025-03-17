package config_test

import (
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestMemory(t *testing.T) {
	m := config.NewMemory(map[string]any{"hello": "world"})
	context := config.GetContext{}
	value, err := m.Get(context)
	if err != nil {
		t.Fatal(err)
	}
	valueMap, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map but found %T", value)
	}
	helloValue, ok := valueMap["hello"]
	if !ok {
		t.Fatal("missing key 'hello'")
	}
	if helloValue != "world" {
		t.Fatalf("expected 'world' found %s", helloValue)
	}
}
