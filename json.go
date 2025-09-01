package config

import (
	"encoding/json"

	"github.com/patrickhuber/go-cross/fs"
)

type JsonOption struct {
	Transformers []Transformer
}

func NewJson(filesystem fs.FS, file string, options ...FileOption) Factory {
	return NewFile(filesystem, file, NewJsonCodec(), options...)
}

type jsonCodec struct{}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}
func (codec *jsonCodec) Unmarshal(buf []byte) (any, error) {
	var data any
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (codec *jsonCodec) Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}
