package config

import (
	"encoding/json"
	"os"
)

type JsonProvider struct {
	file string
}

func NewJson(file string) *JsonProvider {
	return &JsonProvider{
		file: file,
	}
}

func (p *JsonProvider) Get(context GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = json.Unmarshal(buf, &data)
	return data, err
}
