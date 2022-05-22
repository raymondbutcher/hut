package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type terraformCommand struct {
	args []string
	env  []string
}

func NewTerraformCommand(args ...string) *terraformCommand {
	return &terraformCommand{args: args, env: []string{}}
}

func (cmd *terraformCommand) AddArgument(name, value string) {
	arg := name + "=" + value
	if name == "-chdir" {
		cmd.args = append([]string{arg}, cmd.args...)
	} else {
		cmd.args = append(cmd.args, arg)
	}
}

func (cmd *terraformCommand) AddEnvironmentVariable(name, value string) {
	cmd.env = append(cmd.env, name+"="+value)
}

func (cmd terraformCommand) Run() error {
	c := exec.Command("terraform")
	c.Args = cmd.args
	c.Env = cmd.env
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (cmd terraformCommand) String() string {
	parts := []string{}
	parts = append(parts, cmd.env...)
	parts = append(parts, "terraform")
	parts = append(parts, cmd.args...)
	return strings.Join(parts, " ")
}

func (cmd terraformCommand) Subcommand() string {
	for _, arg := range cmd.args {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}
	return ""
}

func (cmd *terraformCommand) Update(module fs.FS, cwd string) error {
	// Skip for most subcommands.
	subcommand := cmd.Subcommand()
	if subcommand != "init" && subcommand != "plan" && subcommand != "apply" {
		return nil
	}

	// Skip if there is no terraform.tfvars[.json] in the current directory.
	if found, err := hasMatchingFiles(module, cwd, "terraform.tfvars", "terraform.tfvars.json"); err != nil {
		return fmt.Errorf("module.HasMatchingFiles: %w", err)
	} else if !found {
		return nil
	}

	// Add the -chdir option to make Terraform change to the module directory when it runs.
	if cwd != "." {
		if chdir, err := filepath.Rel(cwd, "."); err != nil {
			return fmt.Errorf("chdir filepath.Rel: %w", err)
		} else {
			cmd.AddArgument("-chdir", path.Clean(chdir))
		}
	}

	// Changing the directory makes Terraform put the .terraform data directory
	// in that directory, not the current directory. This causes problems when
	// using one module directory and multiple tfbackend files, because Terraform
	// will complain if the backend details don't match those in the .terraform
	// directory. To avoid this conflict, tell Terraform to put the .terraform
	// data directory in the current directory.
	if cwd != "." {
		cmd.AddEnvironmentVariable("TF_DATA_DIR", filepath.Join(cwd, ".terraform"))
	}

	// Automatically use backend files in the current directory and parent directories
	// up to and including the module directory.
	if subcommand == "init" {
		backendFiles, err := findMatchingFilesInTree(module, cwd, "*.tfbackend")
		if err != nil {
			return fmt.Errorf("backendFiles findMatchingFiles: %w", err)
		}
		for _, backendFile := range backendFiles {
			cmd.AddArgument("-backend-config", backendFile)
		}
	}

	// Automatically use var files in the current directory and parent directories
	// up to and including the module directory.
	if subcommand == "plan" || subcommand == "apply" {
		varFiles, err := findMatchingFilesInTree(module, cwd, "terraform.tfvars", "terraform.tfvars.json", "*.auto.tfvars", "*.auto.tfvars.json")
		if err != nil {
			return fmt.Errorf("varFiles findMatchingFiles: %w", err)
		}
		for _, varFile := range varFiles {
			cmd.AddArgument("-var-file", varFile)
		}
	}

	return nil
}
