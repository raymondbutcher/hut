package main

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func findClosestDirectoryWithMatchingFiles(cwd string, patterns ...string) (dir string, err error) {
	if len(patterns) == 0 {
		return "", errors.New("no patterns provided")
	}

	dir = path.Clean(cwd)

	for {
		dirfs := os.DirFS(dir)
		if found, err := hasMatchingFiles(dirfs, ".", patterns...); err != nil {
			return "", fmt.Errorf("hasMatchingFiles: %w", err)
		} else if found {
			return dir, nil
		}

		parent := path.Dir(dir)
		if parent == dir {
			// Reached the top of the filesystem without finding a module directory.
			return "", nil
		}
		dir = parent
	}
}

func findClosestModuleDirectory(cwd string) (dir string, err error) {
	return findClosestDirectoryWithMatchingFiles(cwd, "*.tf", "*.tf.json")
}
