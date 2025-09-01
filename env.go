package config

import (
	"strings"

	"github.com/patrickhuber/go-cross/env"
)

type EnvProvider struct {
	options []EnvOption
	env     env.Environment
}

type EnvOption struct {
	Prefix       string
	Transformers []Transformer
}

func NewEnv(environment env.Environment, options ...EnvOption) Factory {
	provider := &EnvProvider{
		options: options,
		env:     environment,
	}

	return NewFactory(provider)
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
	envVars := p.env.Export()
	for key, value := range envVars {
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
