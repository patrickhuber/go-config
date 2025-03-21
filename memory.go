package config

type memoryProvider struct {
	memory map[string]any
}

func NewMemory(memory map[string]any) Provider {
	return &memoryProvider{
		memory: memory,
	}
}

func (m *memoryProvider) Get(ctx *GetContext) (any, error) {
	return m.memory, nil
}
