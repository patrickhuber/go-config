package config_test

import (
	"log"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

func Example() {
	args := []string{"--hello", "world"}

	// Create a target for cross-platform abstractions
	target := cross.NewTest(platform.Linux, arch.AMD64)
	filesystem := target.FS()

	// Use environment provider from target and set environment variable
	osEnv := target.Env()
	osEnv.Set("env", "yes")

	builder := config.NewBuilder(
		config.NewYaml(filesystem, "config.yml"),
		config.NewJson(filesystem, "config.json"),
		config.NewToml(filesystem, "config.toml"),
		config.NewEnv(osEnv, config.EnvOption{Prefix: "env"}),
		config.NewDotEnv(filesystem, ".env"),
		config.NewFlag([]config.Flag{
			&config.StringFlag{
				Name: "hello",
			},
		}, args),
	)
	root, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := root.Get(&config.GetContext{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%v", cfg)
	}
}
