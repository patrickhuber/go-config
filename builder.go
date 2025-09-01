package config

type Builder struct {
	factories []Factory
}

func NewBuilder(factories ...Factory) *Builder {
	builder := &Builder{
		factories: factories,
	}
	return builder
}

func (b *Builder) WithFactory(factory Factory) *Builder {
	b.factories = append(b.factories, factory)
	return b
}

func (b *Builder) WithProvider(provider Provider) *Builder {
	b.factories = append(b.factories, NewFactory(provider))
	return b
}

func (b *Builder) Build() (Root, error) {
	var providers []Provider
	for _, factory := range b.factories {
		factoryProviders, err := factory.Providers()
		if err != nil {
			return nil, err
		}
		providers = append(providers, factoryProviders...)
	}
	return NewRoot(providers...), nil
}
