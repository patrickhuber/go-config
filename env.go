package config

import (
	"os"
	"strings"
)

type EnvProvider struct {
	prefix     string
	transforms []Transformer
}

func NewEnv(prefix string, transforms ...Transformer) *EnvProvider {
	return &EnvProvider{
		prefix:     prefix,
		transforms: transforms,
	}
}

func (p *EnvProvider) Get(ctx *GetContext) (any, error) {
	prefixSpecified := !strings.EqualFold(p.prefix, "")
	cfg := map[string]any{}

	// load environment variables
	for _, env := range os.Environ() {
		splits := strings.Split(env, "=")
		if len(splits) < 2 {
			continue
		}
		key := splits[0]
		value := splits[1]
		if prefixSpecified && !strings.HasPrefix(key, p.prefix) {
			continue
		}
		cfg[key] = value
	}

	// perform transforms
	return transform(cfg, p.transforms)
}
