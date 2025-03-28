package config

import (
	"os"
)

type fileProvider struct {
	file   string
	codec  Codec
	option FileOption
}

type FileOption struct {
	Transformers []Transformer
}

func NewFile(file string, codec Codec, options ...FileOption) Provider {
	provider := &fileProvider{
		file:  file,
		codec: codec,
	}
	for _, option := range options {
		provider.option.Transformers = append(provider.option.Transformers, option.Transformers...)
	}
	return provider
}

func (provider *fileProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(provider.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = provider.codec.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return transform(data, provider.option.Transformers)
}

func (provider *fileProvider) Set(ctx *SetContext, value any) error {
	buf, err := provider.codec.Marshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(provider.file, buf, 0644)
}
