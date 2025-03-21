package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type tomlProvider struct {
	file         string
	transformers []Transformer
}

func NewToml(file string, transformers ...Transformer) Provider {
	return &tomlProvider{
		file:         file,
		transformers: transformers,
	}
}

func (p *tomlProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = toml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return transform(data, p.transformers)
}
