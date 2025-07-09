package config_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/platform"
)

func TestDynamic(t *testing.T) {
	key := "CONFIG_TEST_FILE_PATH"

	// Use Target for cross-platform abstractions
	target := cross.NewTest(platform.Linux, arch.AMD64)

	// Use the filesystem from the target
	filesystem := target.FS()
	path := target.Path()

	// Use a base directory in the memory filesystem
	testDir := "/test"
	filePath := path.Join(testDir, "abc123.yml")
	fileDirectory := path.Dir(filePath)

	// Ensure directory exists
	exists, err := filesystem.Exists(fileDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		err := filesystem.MkdirAll(fileDirectory, 0666)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = filesystem.WriteFile(filePath, []byte(`key: hello`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a memory-based environment
	envProvider := env.NewMemory()
	envProvider.Set(key, filePath)

	envConfigProvider := config.NewEnv(envProvider, config.EnvOption{Prefix: key})
	dynamicProvider := config.NewDynamic(func(ctx *config.GetContext) (config.Provider, error) {
		m, ok := ctx.MergedConfiguration.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map[string]any but found %T", ctx.MergedConfiguration)
		}
		filePath, ok := m[key].(string)
		if !ok {
			return nil, fmt.Errorf("expected file path to be string but found %T", m[key])
		}
		return config.NewYaml(filesystem, filePath), nil
	})

	root := config.NewRoot(envConfigProvider, dynamicProvider)
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
