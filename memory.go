package config

type MemoryProvider struct {
	memory map[string]any
}

func NewMemory(memory map[string]any) *MemoryProvider {
	return &MemoryProvider{
		memory: memory,
	}
}

func (m *MemoryProvider) Get(context GetContext) (any, error) {
	return m.memory, nil
}
