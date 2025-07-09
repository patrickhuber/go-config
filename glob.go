package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
)

// GlobProviderResolver returns the Provider for the given glob match
type GlobProviderResolver func(filesystem fs.FS, match string) Provider

type GlobOption struct {
	Transformers []Transformer
	Resolver     GlobProviderResolver
}

type globDirection string

const globDirectionUp globDirection = "up"
const globDirectionDown globDirection = "down"

type globProvider struct {
	pattern    string
	directory  string
	direction  globDirection
	filesystem fs.FS
	filepath   filepath.Provider
	options    GlobOption
}

func NewGlob(filesystem fs.FS, filePath filepath.Provider, directory string, pattern string, options ...GlobOption) Provider {
	return newGlobWithDirection(filesystem, filePath, directory, pattern, globDirectionDown, options...)
}

func NewGlobUp(filesystem fs.FS, filePath filepath.Provider, directory string, pattern string, options ...GlobOption) Provider {
	return newGlobWithDirection(filesystem, filePath, directory, pattern, globDirectionUp, options...)
}

func newGlobWithDirection(filesystem fs.FS, filePath filepath.Provider, directory string, pattern string, direction globDirection, options ...GlobOption) Provider {
	provider := &globProvider{
		direction:  direction,
		pattern:    pattern,
		directory:  directory,
		filesystem: filesystem,
		filepath:   filePath,
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
		matches, err = glob(g.filesystem, g.filepath, g.directory, g.pattern)
	case globDirectionUp:
		matches, err = globUp(g.filesystem, g.filepath, g.directory, g.pattern)
	default:
		err = fmt.Errorf("unknow direction")
	}
	if err != nil {
		return nil, err
	}
	var providers []Provider
	for _, match := range matches {
		provider := g.options.Resolver(g.filesystem, match)
		if provider == nil {
			continue
		}
		providers = append(providers, provider)
	}
	root := NewRoot(providers...)

	data, err := root.Get(ctx)
	if err != nil {
		return nil, err
	}
	return transform(data, g.options.Transformers)
}

func defaultGlobProviderResolver(filesystem fs.FS, match string) Provider {
	// Create a filepath provider for this operation
	osProvider := os.New()
	fpProvider := filepath.NewProviderFromOS(osProvider)
	ext := fpProvider.Ext(match)
	var provider Provider
	switch ext {
	case ".json":
		provider = NewJson(filesystem, match)
	case ".yml", ".yaml":
		provider = NewYaml(filesystem, match)
	case ".toml":
		provider = NewToml(filesystem, match)
	case ".env":
		provider = NewDotEnv(filesystem, match)
	}
	return provider
}

func glob(filesystem fs.FS, fp filepath.Provider, dir string, pattern string) ([]string, error) {
	var files []string

	pat := toRegexp(pattern)
	r := regexp.MustCompile(pat)

	// Check if the directory exists in the filesystem
	exists, err := filesystem.Exists(dir)
	if err != nil {
		return nil, err
	}
	if !exists {
		return files, nil
	}

	// Use a simple recursive approach since we can't use fs.WalkDir with memory filesystem
	return globRecursive(filesystem, fp, dir, dir, r)
}

func globRecursive(filesystem fs.FS, fp filepath.Provider, baseDir, currentDir string, r *regexp.Regexp) ([]string, error) {
	var files []string

	entries, err := filesystem.ReadDir(currentDir)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		entryPath := fp.Join(currentDir, entry.Name())

		if entry.IsDir() {
			// Recursively search subdirectories
			subFiles, err := globRecursive(filesystem, fp, baseDir, entryPath, r)
			if err != nil {
				return files, err
			}
			files = append(files, subFiles...)
		} else {
			// Check if file matches pattern
			if r.MatchString(entryPath) {
				files = append(files, entryPath)
			}
		}
	}

	return files, nil
}

func globUp(filesystem fs.FS, fp filepath.Provider, dir string, pattern string) ([]string, error) {
	var files []string
	current := fp.Clean(dir)

	pat := toRegexp(pattern)
	r := regexp.MustCompile(pat)

	for {
		parentDir := fp.Dir(current)
		if strings.Compare(current, parentDir) == 0 {
			break
		}

		// Check if current directory exists in the filesystem
		exists, err := filesystem.Exists(current)
		if err != nil {
			return nil, err
		}
		if !exists {
			current = parentDir
			continue
		}

		entries, err := filesystem.ReadDir(current)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			path := fp.Join(current, name)
			if r.MatchString(path) {
				files = append(files, path)
			}
		}
		current = parentDir
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
