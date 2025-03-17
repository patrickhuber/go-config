package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TomlProvider struct {
	file string
}

func NewToml(file string) *TomlProvider {
	return &TomlProvider{
		file: file,
	}
}

func (p *TomlProvider) Get(context GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = toml.Unmarshal(buf, &data)
	return data, err
}
