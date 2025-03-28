package config

import "gopkg.in/yaml.v3"

func NewYaml(file string, options ...FileOption) Provider {
	return NewFile(file, NewYamlCodec(), options...)
}

type yamlCodec struct{}

func NewYamlCodec() Codec {
	return &yamlCodec{}
}

func (codec *yamlCodec) Unmarshal(buf []byte, data any) error {
	return yaml.Unmarshal(buf, data)
}

func (codec *yamlCodec) Marshal(data any) ([]byte, error) {
	return yaml.Marshal(data)
}
