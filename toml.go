package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type tomlProvider struct {
	file string
}

func NewToml(file string) Provider {
	return &tomlProvider{
		file: file,
	}
}

func (p *tomlProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = toml.Unmarshal(buf, &data)
	return data, err
}
