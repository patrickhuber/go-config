package config_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-config"
)

func TestBuilder(t *testing.T) {
	expected := map[string]any{
		"hello": "world",
	}
	builder := config.NewBuilder(&FakeProvider{
		Data: expected,
	})
	actual, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	actualMap, ok := actual.(map[string]any)
	if !ok {
		t.Fatalf("expected map but found %T", actual)
	}
	for expectedKey, expectedValue := range expected {
		actualValue, ok := actualMap[expectedKey]
		if !ok {
			t.Fatalf("expected to find key %s in result", expectedKey)
		}
		if !reflect.DeepEqual(expectedValue, actualValue) {
			t.Fatalf("expected %v to equal %v", expectedValue, actualValue)
		}
	}
	for actualKey, actualValue := range actualMap {
		_, ok := expected[actualKey]
		if !ok {
			t.Fatalf("found unexpected key %s and value %s in result", actualKey, actualValue)
		}
	}
}

type FakeProvider struct {
	Data  any
	Error error
}

func (p *FakeProvider) Get(context config.GetContext) (any, error) {
	if p.Error != nil {
		return nil, p.Error
	}
	return p.Data, nil
}
