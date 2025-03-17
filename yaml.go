package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YamlProvider struct {
	file string
}

func NewYaml(file string) *YamlProvider {
	return &YamlProvider{
		file: file,
	}
}

func (p *YamlProvider) Get() (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = yaml.Unmarshal(buf, &data)
	return data, err
}
