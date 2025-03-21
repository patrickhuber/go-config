package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type yamlProvider struct {
	file         string
	transformers []Transformer
}

func NewYaml(file string, transfomers ...Transformer) Provider {
	return &yamlProvider{
		file:         file,
		transformers: transfomers,
	}
}

func (p *yamlProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return transform(data, p.transformers)
}
