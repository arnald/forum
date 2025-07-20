package path

import (
	"os"
	"path/filepath"
)

type Resolver struct {
	projectRoot string
}

func NewResolver() *Resolver {
	cwd, _ := os.Getwd()
	current := cwd
	for {
		_, err := os.Stat(filepath.Join(current, "go.mod"))
		if err == nil {
			return &Resolver{projectRoot: current}
		}

		parent := filepath.Dir(current)
		if parent == current {
			return &Resolver{projectRoot: cwd}
		}
		current = parent
	}
}

func (r *Resolver) GetPath(relPath string) string {
	return filepath.Join(r.projectRoot, filepath.FromSlash(relPath))
}
