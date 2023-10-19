package resource

import (
	"os"
	"path"
	"path/filepath"
)

// Scanner scans files using the informed glob pattern, returns a slice with the the files
// matching the expression.
func Scanner(pattern string) ([]string, error) {
	patterns := []string{}

	// when a directory is informed the patterns will only select ".yaml" and ".yml" files
	info, _ := os.Stat(pattern)
	if info != nil && info.IsDir() {
		patterns = append(patterns, path.Join(pattern, "*.yml"))
		patterns = append(patterns, path.Join(pattern, "*.yaml"))
	}

	files := []string{}
	for _, p := range patterns {
		match, err := filepath.Glob(p)
		if err != nil {
			return nil, err
		}
		files = append(files, match...)
	}
	return files, nil
}
