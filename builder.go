package config

type Builder struct {
	providers []Provider
}

func NewBuilder(providers ...Provider) *Builder {
	builder := &Builder{
		providers: providers,
	}
	return builder
}

func (b *Builder) With(provider Provider) *Builder {
	b.providers = append(b.providers, provider)
	return b
}

func (b *Builder) Build() (Root, error) {
	return NewRoot(b.providers...), nil
}
