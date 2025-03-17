package config

type Provider interface {
	Get(context GetContext) (any, error)
}

type GetContext struct {
}
