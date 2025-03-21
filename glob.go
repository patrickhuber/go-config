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

func NewGlob(directory string, pattern string, transformers ...Transformer) Provider {
	return NewGlobWithResolver(directory, pattern, defaultGlobProviderResolver, transformers...)
}

func NewGlobWithResolver(directory string, pattern string, resolver GlobProviderResolver, transformers ...Transformer) Provider {
	return &globProvider{
		direction:    globDirectionDown,
		pattern:      pattern,
		directory:    directory,
		transformers: transformers,
		resolver:     resolver,
	}
}

func NewGlobUp(directory string, pattern string, transformers ...Transformer) Provider {
	return NewGlobUpWithResolver(directory, pattern, defaultGlobProviderResolver, transformers...)
}

func NewGlobUpWithResolver(directory string, pattern string, resolver GlobProviderResolver, transformers ...Transformer) Provider {
	return &globProvider{
		direction:    globDirectionUp,
		pattern:      pattern,
		directory:    directory,
		transformers: transformers,
		resolver:     resolver,
	}
}

type globDirection string

const globDirectionUp globDirection = "up"
const globDirectionDown globDirection = "down"

type globProvider struct {
	pattern      string
	directory    string
	direction    globDirection
	resolver     GlobProviderResolver
	transformers []Transformer
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
		provider := g.resolver(match)
		if provider == nil {
			continue
		}
		providers = append(providers, provider)
	}
	build, err := NewBuilder(providers...).Build()
	if err != nil {
		return nil, err
	}
	return transform(build, g.transformers)
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
