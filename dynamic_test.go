package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
)

func TestDynamic(t *testing.T) {
	key := "CONFIG_TEST_FILE_PATH"
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "abc123.yml")
	err := os.WriteFile(filePath, []byte(`key: hello`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv(key, filePath)
	osEnv := env.New()
	envProvider := config.NewEnv(osEnv, config.EnvOption{Prefix: key})
	dynamicProvider := config.NewDynamic(func(ctx *config.GetContext) (config.Provider, error) {
		m, ok := ctx.MergedConfiguration.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map[string]any but found %T", ctx.MergedConfiguration)
		}
		filePath, ok := m[key].(string)
		if !ok {
			return nil, fmt.Errorf("expected file path to be string but found %T", m[key])
		}
		return config.NewYaml(filePath), nil
	})

	root := config.NewRoot(envProvider, dynamicProvider)
	cfg, err := root.Get(&config.GetContext{})
	if err != nil {
		t.Fatal(err)
	}

	cfgMap, ok := cfg.(map[string]any)
	if !ok {
		t.Fatalf("expected config to be map but found %T", cfg)
	}
	value := cfgMap["key"]
	if value != "hello" {
		t.Fatalf("expected cfgMap['key'] to be 'hello' but it was %v", value)
	}
}
