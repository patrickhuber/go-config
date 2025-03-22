package config

import (
	"github.com/BurntSushi/toml"
)

func NewToml(file string, transformers ...Transformer) Provider {
	return NewFile(file, NewTomlCodec(), transformers...)
}

type tomlCodec struct{}

func NewTomlCodec() Codec {
	return &tomlCodec{}
}
func (codec *tomlCodec) Unmarshal(buf []byte, data any) error {
	return toml.Unmarshal(buf, data)
}

func (codec *tomlCodec) Marshal(data any) ([]byte, error) {
	return toml.Marshal(data)
}
