package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cmd := NewTerraformCommand(os.Args[1:]...)

	cwd, err := os.Getwd()
	exitErr(err)

	moduledir, err := findClosestModuleDirectory(cwd)
	exitErr(err)

	if moduledir != "" {
		cwd, err := filepath.Rel(moduledir, cwd)
		exitErr(err)

		err = cmd.Update(os.DirFS(moduledir), cwd)
		exitErr(err)
	}

	if os.Getenv("HUT_DRY_RUN") != "" {
		fmt.Fprintf(os.Stderr, "# hut dry run: %s\n", cmd)
	} else {
		fmt.Fprintf(os.Stderr, "# hut run: %s\n", cmd)
		err = cmd.Run()
		exitErr(err)
	}
}

func exitErr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "# hut error: %s\n", err)

	if exitErr, isExitError := err.(*exec.ExitError); isExitError {
		os.Exit(exitErr.ExitCode())
	} else {
		os.Exit(1)
	}
}
