package main

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

func findMatchingFilesInTree(fsys fs.FS, cwd string, patterns ...string) ([]string, error) {
	files := []string{}

	if len(patterns) == 0 {
		return files, errors.New("no patterns provided")
	}

	dirs := []string{}
	parts := []string{}
	for _, part := range strings.Split(cwd, "/") {
		parts = append(parts, part)
		dirs = append(dirs, strings.Join(parts, "/"))
	}

	for _, dir := range dirs {
		ents, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return files, fmt.Errorf("fs.ReadDir: %w", err)
		}
		for _, pattern := range patterns {
			for _, ent := range ents {
				if !ent.IsDir() {
					match, err := path.Match(pattern, ent.Name())
					if err != nil {
						return files, fmt.Errorf("path.Match: %w", err)
					}
					if match {
						files = append(files, dir+"/"+ent.Name())
					}
				}
			}
		}
	}
	return files, nil
}

func hasMatchingFiles(fsys fs.FS, dir string, patterns ...string) (bool, error) {
	if len(patterns) == 0 {
		return false, errors.New("no patterns provided")
	}

	ents, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return false, fmt.Errorf("fs.ReadDir: %w", err)
	}
	for _, ent := range ents {
		if ent.IsDir() {
			continue
		}
		for _, pattern := range patterns {
			match, err := path.Match(pattern, ent.Name())
			if err != nil {
				return false, fmt.Errorf("path.Match: %w", err)
			}
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}
