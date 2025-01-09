package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type GlobUpProvider struct {
	directory string
	pattern   string
}

func NewGlobUp(directory string, pattern string) *GlobUpProvider {
	return &GlobUpProvider{
		directory: directory,
		pattern:   pattern,
	}
}

func (p *GlobUpProvider) Get() (any, error) {
	matches, err := globUp(p.directory, p.pattern)
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
			if r.MatchString(name) {
				path := filepath.Join(current, name)
				files = append(files, path)
			}
		}
		current = dir
	}
	return files, nil
}
