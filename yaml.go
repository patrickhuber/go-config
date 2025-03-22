package config

import "gopkg.in/yaml.v3"

func NewYaml(file string, transfomers ...Transformer) Provider {
	return NewFile(file, NewYamlCodec(), transfomers...)
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
