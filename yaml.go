package config

import (
	"gopkg.in/yaml.v3"

	"github.com/patrickhuber/go-cross/fs"
)

func NewYaml(filesystem fs.FS, file string, options ...FileOption) Provider {
	return NewFile(filesystem, file, NewYamlCodec(), options...)
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
