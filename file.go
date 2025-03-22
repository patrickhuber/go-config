package config

import (
	"os"
)

type fileProvider struct {
	file         string
	codec        Codec
	transformers []Transformer
}

func NewFile(file string, codec Codec, transformers ...Transformer) Provider {
	return &fileProvider{
		file:         file,
		codec:        codec,
		transformers: transformers,
	}
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
	return transform(data, provider.transformers)
}
