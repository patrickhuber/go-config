package config

type Builder struct {
	providers []Provider
}

func NewBuilder(providers ...Provider) *Builder {
	builder := &Builder{}
	for _, provider := range providers {
		builder.With(provider)
	}
	return builder
}

func (b *Builder) With(provider Provider) *Builder {
	b.providers = append(b.providers, provider)
	return b
}

func (b *Builder) Build() (any, error) {
	var result any = nil
	ctx := &GetContext{}
	for _, provider := range b.providers {
		cfg, err := provider.Get(ctx)
		if err != nil {
			return nil, err
		}
		if result == nil {
			result = cfg
			ctx.MergedConfiguration = cfg
			continue
		}
		result, err = Merge(cfg, result)
		if err != nil {
			return nil, err
		}
		ctx.MergedConfiguration = result
	}
	return result, nil
}
