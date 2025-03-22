package config

import (
	"encoding/json"
)

func NewJson(file string, transformers ...Transformer) Provider {
	return NewFile(file, NewJsonCodec(), transformers...)
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
