package config

import (
	"github.com/patrickhuber/go-cross/fs"
)

type fileProvider struct {
	file    string
	codec   Codec
	fs      fs.FS
	options []FileOption
}

type FileBeforeGet func(ctx *GetContext, provider Provider, file string) (any, error)

type FileOption struct {
	Transformers []Transformer
	BeforeGet    FileBeforeGet
}

func NewFile(filesystem fs.FS, file string, codec Codec, options ...FileOption) Provider {
	provider := &fileProvider{
		file:    file,
		codec:   codec,
		fs:      filesystem,
		options: options,
	}
	return provider
}

func (provider *fileProvider) Get(ctx *GetContext) (any, error) {
	for _, option := range provider.options {
		if option.BeforeGet == nil {
			continue
		}
		beforeGet, err := option.BeforeGet(ctx, provider, provider.file)
		if err != nil {
			return nil, err
		}
		ctx.MergedConfiguration, err = Merge(ctx.MergedConfiguration, beforeGet)
		if err != nil {
			return nil, err
		}
	}
	buf, err := provider.fs.ReadFile(provider.file)
	if err != nil {
		return nil, err
	}
	data, err := provider.codec.Unmarshal(buf)
	if err != nil {
		return nil, err
	}
	var transformers []Transformer
	for _, option := range provider.options {
		transformers = append(transformers, option.Transformers...)
	}
	return transform(data, transformers)
}

func (provider *fileProvider) Set(ctx *SetContext, value any) error {
	buf, err := provider.codec.Marshal(value)
	if err != nil {
		return err
	}
	return provider.fs.WriteFile(provider.file, buf, 0644)
}
