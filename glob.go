package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
)

// GlobResolver returns the Provider for the given glob match
type GlobResolver interface {
	Resolve(match string) Factory
}

func DefaultGlobResolver(fileSystem fs.FS, path filepath.Provider) GlobResolver {
	return &defaultGlobResolver{
		fs:   fileSystem,
		path: path,
	}
}

type defaultGlobResolver struct {
	fs   fs.FS
	path filepath.Provider
}

func (r *defaultGlobResolver) Resolve(match string) Factory {
	var factory Factory
	ext := r.path.Ext(match)
	switch ext {
	case ".json":
		factory = NewJson(r.fs, match)
	case ".yml", ".yaml":
		factory = NewYaml(r.fs, match)
	case ".toml":
		factory = NewToml(r.fs, match)
	case ".env":
		factory = NewDotEnv(r.fs, match)
	}
	return factory
}

type GlobOption struct {
	Transformers []Transformer
}

type globDirection string

const globDirectionUp globDirection = "up"
const globDirectionDown globDirection = "down"

type globProviderFactory struct {
	pattern    string
	directory  string
	direction  globDirection
	filesystem fs.FS
	filepath   filepath.Provider
	resolver   GlobResolver
	options    []GlobOption
}

func NewGlob(
	filesystem fs.FS,
	filePath filepath.Provider,
	resolver GlobResolver,
	directory string,
	pattern string,
	options ...GlobOption) Factory {
	return newGlobProviderFactoryWithDirection(
		filesystem,
		filePath,
		resolver,
		globDirectionDown,
		directory,
		pattern,
		options...)
}

func NewGlobUp(
	filesystem fs.FS,
	filePath filepath.Provider,
	resolver GlobResolver,
	directory string,
	pattern string,
	options ...GlobOption) Factory {
	return newGlobProviderFactoryWithDirection(
		filesystem,
		filePath,
		resolver,
		globDirectionUp,
		directory,
		pattern,
		options...)
}

func newGlobProviderFactoryWithDirection(
	filesystem fs.FS,
	filePath filepath.Provider,
	resolver GlobResolver,
	direction globDirection,
	directory string,
	pattern string,
	options ...GlobOption) Factory {
	return &globProviderFactory{
		pattern:    pattern,
		directory:  directory,
		direction:  direction,
		filesystem: filesystem,
		filepath:   filePath,
		resolver:   resolver,
		options:    options,
	}
}

func (g *globProviderFactory) Providers() ([]Provider, error) {
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
		factory := g.resolver.Resolve(match)
		if factory == nil {
			continue
		}
		childProviders, err := factory.Providers()
		if err != nil {
			return nil, err
		}
		providers = append(providers, childProviders...)
	}
	return providers, nil
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
