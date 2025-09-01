package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/patrickhuber/go-cross/fs"
)

func NewDotEnv(filesystem fs.FS, file string, options ...FileOption) Factory {
	return NewFile(filesystem, file, &dotEnvCodec{}, options...)
}

type dotEnvCodec struct{}

func (codec *dotEnvCodec) Marshal(data any) ([]byte, error) {
	return nil, fmt.Errorf("dotEnvCodec Marshal is not implemented")
}

func (codec *dotEnvCodec) Unmarshal(buf []byte) (any, error) {
	stringAnyMap := make(map[string]any)
	stringStringMap, err := godotenv.UnmarshalBytes(buf)
	if err != nil {
		return nil, err
	}
	for k, v := range stringStringMap {
		stringAnyMap[k] = v
	}
	return stringAnyMap, nil
}
