package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DotEnvProvider struct {
	file string
}

func NewDotEnv(file string) *DotEnvProvider {
	return &DotEnvProvider{
		file: file,
	}
}

func (p *DotEnvProvider) Get(ctx *GetContext) (any, error) {
	file, err := os.Open(p.file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]any)
	kv, err := godotenv.Parse(file)
	if err != nil {
		return nil, err
	}
	for k, v := range kv {
		result[k] = v
	}
	return result, nil
}
