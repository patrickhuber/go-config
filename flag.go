package config

import (
	"flag"
	"fmt"
)

type Flag interface {
	Set(value string) error
	String() string
	Value() any
}

type StringFlag struct {
	Name    string
	Default string
	Usage   string
	value   string
}

func (s *StringFlag) Set(value string) error {
	s.value = value
	return nil
}

func (s *StringFlag) String() string {
	return s.value
}

func (s *StringFlag) Value() any {
	return s.value
}

type StringSliceFlag struct {
	Name    string
	Default []string
	Usage   string
	value   []any
}

func (s *StringSliceFlag) Set(value string) error {
	s.value = append(s.value, value)
	return nil
}

func (s *StringSliceFlag) String() string {
	return fmt.Sprintf("%v", s.value)
}

func (s *StringSliceFlag) Value() any {
	return s.value
}

type flagProvider struct {
	flags   []Flag
	args    []string
	options []FlagOption
}

type FlagOption struct {
	Transformers []Transformer
}

func NewFlag(flags []Flag, args []string, options ...FlagOption) Factory {
	return NewFactory(
		&flagProvider{
			flags:   flags,
			args:    args,
			options: options,
		},
	)
}

func (p *flagProvider) Get(ctx *GetContext) (any, error) {
	m := map[string]any{}
	flagset := flag.NewFlagSet("any", flag.ContinueOnError)
	for _, f := range p.flags {
		switch ty := f.(type) {
		case *StringFlag:
			flagset.Var(ty, ty.Name, ty.Usage)
		case *StringSliceFlag:
			flagset.Var(ty, ty.Name, ty.Usage)
		}
	}

	// parse flags
	err := flagset.Parse(p.args)
	if err != nil {
		return nil, err
	}

	// visit all flags and set the flags in the map
	flagset.Visit(func(f *flag.Flag) {
		flagValue, ok := f.Value.(Flag)
		if !ok {
			return
		}
		m[f.Name] = flagValue.Value()
	})

	// run the transforms
	var result any = m
	for _, option := range p.options {
		for _, transformer := range option.Transformers {
			result, err = transformer.Transform(result)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
