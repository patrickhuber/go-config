package config

type memoryProvider struct {
	memory       any
	transformers []Transformer
}

func NewMemory(memory any, transformers ...Transformer) Provider {
	return &memoryProvider{
		memory:       memory,
		transformers: transformers,
	}
}

func (m *memoryProvider) Get(ctx *GetContext) (any, error) {
	return transform(m.memory, m.transformers)
}
