package config

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
