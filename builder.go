package config

import "fmt"

type Builder struct {
	providers []Provider
}

func NewBuilder(providers ...Provider) *Builder {
	builder := &Builder{}
	for _, provider := range providers {
		builder.With(provider)
	}
	return builder
}

func (b *Builder) With(provider Provider) *Builder {
	b.providers = append(b.providers, provider)
	return b
}

func (b *Builder) Build() (any, error) {
	var result any = nil
	for _, provider := range b.providers {
		cfg, err := provider.Get()
		if err != nil {
			return nil, err
		}
		if result == nil {
			result = cfg
			continue
		}
		result, err = merge(cfg, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func merge(source any, target any) (any, error) {
	sourceMap, ok := source.(map[string]any)
	if ok {
		targetMap, ok := target.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("source is a %T and target is a %T", source, target)
		}
		for k, v := range sourceMap {
			targetMap[k] = v
		}
	}
	return target, nil
}
