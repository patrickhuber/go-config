package config

import (
	"encoding/json"
	"os"
)

type jsonProvider struct {
	file string
}

func NewJson(file string) Provider {
	return &jsonProvider{
		file: file,
	}
}

func (p *jsonProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = json.Unmarshal(buf, &data)
	return data, err
}
