package config

type Provider interface {
	Get() (any, error)
}
