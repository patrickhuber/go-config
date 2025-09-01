package config

type Factory interface {
	Providers() ([]Provider, error)
}

type providerFactory struct {
	providers []Provider
}

func NewFactory(providers ...Provider) Factory {
	return &providerFactory{
		providers: providers,
	}
}

func (f *providerFactory) Providers() ([]Provider, error) {
	return f.providers, nil
}
