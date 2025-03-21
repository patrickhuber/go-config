package config

import (
	"encoding/json"
	"os"
)

type jsonProvider struct {
	file         string
	transformers []Transformer
}

func NewJson(file string, transformers ...Transformer) Provider {
	return &jsonProvider{
		file:         file,
		transformers: transformers,
	}
}

func (p *jsonProvider) Get(ctx *GetContext) (any, error) {
	buf, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	var data any
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return transform(data, p.transformers)
}
