package main

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestFindClosestModuleDirectory(t *testing.T) {
	is := is.New(t)

	// Assume tests are running from the repo root.
	root, err := os.Getwd()
	is.NoErr(err)

	// Call findClosestModuleDirectory with a few subdirectories.
	subdirs := []string{
		"example",
		"example/au",
		"example/au/dev",
	}
	for _, subdir := range subdirs {
		t.Run(subdir, func(t *testing.T) {
			// It should always return the example directory
			// because it's the only directory in the tree with a matching file (main.tf).
			is := is.New(t)
			returneddir, err := findClosestModuleDirectory(root + "/" + subdir)
			is.NoErr(err)
			is.Equal(returneddir, root+"/example")
		})
	}
}
