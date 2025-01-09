package config

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

func NewGlob(directory string, pattern string) *GlobProvider {
	return &GlobProvider{
		pattern:   pattern,
		directory: directory,
	}
}

type GlobProvider struct {
	pattern   string
	directory string
}

func (g *GlobProvider) Get() (any, error) {
	matches, err := glob(g.directory, g.pattern)
	if err != nil {
		return nil, err
	}
	var providers []Provider
	for _, match := range matches {
		ext := filepath.Ext(match)
		var provider Provider
		provider = nil
		switch ext {
		case ".json":
			provider = NewJson(match)
		case ".yml", ".yaml":
			provider = NewYaml(match)
		case ".toml":
			provider = NewToml(match)
		default:
			continue
		}
		if provider == nil {
			continue
		}
		providers = append(providers, provider)
	}
	return NewBuilder(providers...).Build()
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
