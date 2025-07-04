package config

type Root interface {
	Provider
	Providers() []Provider
}

type root struct {
	providers []Provider
}

func NewRoot(providers ...Provider) Root {
	return &root{
		providers: providers,
	}
}

func (r *root) Providers() []Provider {
	return r.providers
}

func (r *root) Get(ctx *GetContext) (any, error) {
	current := ctx.MergedConfiguration
	for _, provider := range r.providers {
		currentCtx := &GetContext{
			MergedConfiguration: current,
		}
		cfg, err := provider.Get(currentCtx)
		if err != nil {
			return nil, err
		}
		current, err = Merge(cfg, current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
