package config

import "gopkg.in/yaml.v3"

func NewYaml(file string, options ...FileOption) Provider {
	return NewFile(file, NewYamlCodec(), options...)
}

type yamlCodec struct{}

func NewYamlCodec() Codec {
	return &yamlCodec{}
}

func (codec *yamlCodec) Unmarshal(buf []byte) (any, error) {
	var data any
	err := yaml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (codec *yamlCodec) Marshal(data any) ([]byte, error) {
	return yaml.Marshal(data)
}
