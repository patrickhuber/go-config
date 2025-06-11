package config

import "fmt"

type Transformer interface {
	Transform(instance any) (any, error)
}

type funcTransformer struct {
	transform func(any) (any, error)
}

func (t *funcTransformer) Transform(instance any) (any, error) {
	return t.transform(instance)
}

func FuncTransformer(transform func(any) (any, error)) Transformer {
	return &funcTransformer{
		transform: transform,
	}
}

func TransformProvider(transform func(any) (any, error)) Provider {
	return &transformerProvider{
		transformer: FuncTransformer(transform),
	}
}

func FuncTypedTransformer[T any](transform func(T) (T, error)) Transformer {
	return &funcTransformer{
		transform: func(a any) (any, error) {

			// if config is nil, do no transformation
			if a == nil {
				return a, nil
			}

			// if config is not a map, throw an error
			m, ok := a.(T)
			if !ok {
				var zero T
				return nil, fmt.Errorf("unable to convert %T to %T", a, zero)
			}

			// run the transform function over the map
			return transform(m)
		},
	}
}

type transformerProvider struct {
	transformer Transformer
}

func (p *transformerProvider) Get(ctx *GetContext) (any, error) {
	return p.transformer.Transform(ctx.MergedConfiguration)
}

func transform(cfg any, transformers []Transformer) (any, error) {
	var err error
	var current any = cfg
	for _, transform := range transformers {
		current, err = transform.Transform(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
