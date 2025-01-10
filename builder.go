package config

import (
	"fmt"
	"reflect"
)

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
	// check if the types are equal, error if not
	if reflect.TypeOf(source) != reflect.TypeOf(target) {
		return nil, fmt.Errorf("%T expected to be equal to %T", source, target)
	}

	switch src := source.(type) {
	case map[string]any:
		// we have a type guard above
		targetMap := target.(map[string]any)
		for k, v := range src {
			targetMap[k] = v
		}
		return targetMap, nil
	}
	return nil, fmt.Errorf("unable to merge type %T", source)
}
