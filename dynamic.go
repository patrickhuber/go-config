package config

type DynamicResolver func(ctx *GetContext) (Provider, error)

type dynamic struct {
	resolver DynamicResolver
	provider Provider
}

func NewDynamic(resolver DynamicResolver) Provider {
	return &dynamic{
		resolver: resolver,
	}
}

// Get implements Provider.
func (d *dynamic) Get(ctx *GetContext) (any, error) {
	provider, err := d.getProvider(ctx)
	if err != nil {
		return nil, err
	}
	return provider.Get(ctx)
}

func (d *dynamic) getProvider(ctx *GetContext) (Provider, error) {
	if d.provider != nil {
		return d.provider, nil
	}

	provider, err := d.resolver(ctx)
	if err != nil {
		return nil, err
	}
	d.provider = provider

	return d.provider, nil
}
