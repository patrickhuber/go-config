package config

type Root interface {
	Provider
	Providers() ([]Provider, error)
}

type root struct {
	providers []Provider
}

func NewRoot(providers ...Provider) Root {
	return &root{
		providers: providers,
	}
}

func (r *root) Providers() ([]Provider, error) {
	return r.providers, nil
}

func (r *root) Get(ctx *GetContext) (any, error) {
	current := ctx.MergedConfiguration
	providers, err := r.Providers()
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
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
