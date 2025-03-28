package config

import (
	"github.com/BurntSushi/toml"
)

func NewToml(file string, options ...FileOption) Provider {
	return NewFile(file, NewTomlCodec(), options...)
}

type tomlCodec struct{}

func NewTomlCodec() Codec {
	return &tomlCodec{}
}
func (codec *tomlCodec) Unmarshal(buf []byte) (any, error) {
	var data any
	err := toml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (codec *tomlCodec) Marshal(data any) ([]byte, error) {
	return toml.Marshal(data)
}
