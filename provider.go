package config

type Provider interface {
	Get(ctx *GetContext) (any, error)
}

type GetContext struct {
	MergedConfiguration any
}

type SetContext struct{}
