package config

import (
	"os"
	"strings"
)

type EnvProvider struct {
	prefix string
}

func NewEnv(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix: prefix,
	}
}

func (p *EnvProvider) Get() (any, error) {
	prefixSpecified := !strings.EqualFold(p.prefix, "")
	cfg := map[string]any{}
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
	return cfg, nil
}
