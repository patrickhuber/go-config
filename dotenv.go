package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func NewDotEnv(file string, options ...FileOption) Provider {
	return NewFile(file, &dotEnvCodec{}, options...)
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
