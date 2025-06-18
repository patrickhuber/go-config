package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GlobProviderResolver returns the Provider for the given glob match
type GlobProviderResolver func(match string) Provider

type GlobOption struct {
	Transformers []Transformer
	Resolver     GlobProviderResolver
}

type globDirection string

const globDirectionUp globDirection = "up"
const globDirectionDown globDirection = "down"

type globProvider struct {
	pattern   string
	directory string
	direction globDirection
	options   GlobOption
}

func NewGlob(directory string, pattern string, options ...GlobOption) Provider {
	return newGlobWithDirection(directory, pattern, globDirectionDown, options...)
}

func NewGlobUp(directory string, pattern string, options ...GlobOption) Provider {
	return newGlobWithDirection(directory, pattern, globDirectionUp, options...)
}

func newGlobWithDirection(directory string, pattern string, direction globDirection, options ...GlobOption) Provider {
	provider := &globProvider{
		direction: direction,
		pattern:   pattern,
		directory: directory,
	}
	for _, option := range options {
		provider.options.Transformers = append(provider.options.Transformers, option.Transformers...)
		provider.options.Resolver = option.Resolver
	}
	if provider.options.Resolver == nil {
		provider.options.Resolver = defaultGlobProviderResolver
	}
	return provider
}

func (g *globProvider) Get(ctx *GetContext) (any, error) {
	var matches []string
	var err error
	switch g.direction {
	case globDirectionDown:
		matches, err = glob(g.directory, g.pattern)
	case globDirectionUp:
		matches, err = globUp(g.directory, g.pattern)
	default:
		err = fmt.Errorf("unknow direction")
	}
	if err != nil {
		return nil, err
	}
	var providers []Provider
	for _, match := range matches {
		provider := g.options.Resolver(match)
		if provider == nil {
			continue
		}
		providers = append(providers, provider)
	}
	root, err := NewBuilder(providers...).Build()
	if err != nil {
		return nil, err
	}
	build, err := root.Get(ctx)
	if err != nil {
		return nil, err
	}
	return transform(build, g.options.Transformers)
}

func defaultGlobProviderResolver(match string) Provider {
	ext := filepath.Ext(match)
	var provider Provider
	switch ext {
	case ".json":
		provider = NewJson(match)
	case ".yml", ".yaml":
		provider = NewYaml(match)
	case ".toml":
		provider = NewToml(match)
	case ".env":
		provider = NewDotEnv(match)
	}
	return provider
}

func glob(dir string, pattern string) ([]string, error) {
	var files []string

	pat := toRegexp(pattern)
	r := regexp.MustCompile(pat)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() || err != nil {
			return nil
		}
		if r.MatchString(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func globUp(dir string, pattern string) ([]string, error) {
	var files []string
	current := filepath.Clean(dir)

	pat := toRegexp(pattern)
	r := regexp.MustCompile(pat)

	for {
		dir := filepath.Dir(current)
		if strings.Compare(current, dir) == 0 {
			break
		}
		entries, err := os.ReadDir(current)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			path := filepath.Join(current, name)
			if r.MatchString(path) {
				files = append(files, path)
			}
		}
		current = dir
	}
	return files, nil
}

var replaces = regexp.MustCompile(`(\.)|(\*\*\/)|(\*)|([^\/\*]+)|(\/)`)

func toRegexp(pattern string) string {
	pat := replaces.ReplaceAllStringFunc(pattern, func(s string) string {
		switch s {
		case "/":
			return "\\/"
		case ".":
			return "\\."
		case "**/":
			return ".*"
		case "*":
			return "[^/]*"
		default:
			return s
		}
	})
	return "^" + pat + "$"
}
