package config

import (
	"encoding/json"
)

type JsonOption struct {
	Transformers []Transformer
}

func NewJson(file string, options ...FileOption) Provider {
	return NewFile(file, NewJsonCodec(), options...)
}

type jsonCodec struct{}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}
func (codec *jsonCodec) Unmarshal(buf []byte, data any) error {
	return json.Unmarshal(buf, data)
}

func (codec *jsonCodec) Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}
