package config

type Marshaler interface {
	Marshal(data any) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal(buf []byte, data any) error
}

type Codec interface {
	Marshaler
	Unmarshaler
}
