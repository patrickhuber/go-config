package config

import (
	"os"
	"strings"
)

type EnvProvider struct {
	options []EnvOption
}

type EnvOption struct {
	Prefix       string
	Transformers []Transformer
}

func NewEnv(options ...EnvOption) *EnvProvider {
	return &EnvProvider{
		options: options,
	}
}

func (p *EnvProvider) Get(ctx *GetContext) (any, error) {
	prefix := ""
	prefixSpecified := false
	for _, option := range p.options {
		if strings.EqualFold(option.Prefix, "") {
			continue
		}
		prefixSpecified = true
		prefix = option.Prefix
	}
	cfg := map[string]any{}

	// load environment variables
	for _, env := range os.Environ() {
		splits := strings.Split(env, "=")
		if len(splits) < 2 {
			continue
		}
		key := splits[0]
		value := splits[1]
		if prefixSpecified && !strings.HasPrefix(key, prefix) {
			continue
		}
		cfg[key] = value
	}

	var transformers []Transformer
	for _, option := range p.options {
		transformers = append(transformers, option.Transformers...)
	}
	// perform transforms
	return transform(cfg, transformers)
}
