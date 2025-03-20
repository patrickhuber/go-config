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
