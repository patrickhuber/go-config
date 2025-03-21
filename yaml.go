package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type yamlProvider struct {
	file string
}

func NewYaml(file string) Provider {
	return &yamlProvider{
		file: file,
	}
}

func (p *yamlProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = yaml.Unmarshal(buf, &data)
	return data, err
}
