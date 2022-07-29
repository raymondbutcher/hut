package main

import (
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
)

func TestTerraformCommandUpdate(t *testing.T) {
	is := is.New(t)

	module := fstest.MapFS{
		"main.tf":                    {},
		"eu/eu.auto.tfvars":          {},
		"eu/eu.s3.tfbackend":         {},
		"eu/dev/terraform.tfvars":    {},
		"eu/dev/terraform.tfbackend": {},
		"us/us.auto.tfvars":          {},
		"us/us.s3.tfbackend":         {},
		"us/dev/terraform.tfvars":    {},
		"us/dev/terraform.tfbackend": {},
	}
	cwd := "eu/dev"

	cmd := NewTerraformCommand("fmt")
	err := cmd.Update(module, cwd)
	is.NoErr(err)
	is.Equal(cmd.String(), "terraform fmt")

	cmd = NewTerraformCommand("init")
	err = cmd.Update(module, cwd)
	is.NoErr(err)
	is.Equal(cmd.String(), "TF_DATA_DIR=eu/dev/.terraform terraform -chdir=../.. init -backend-config=eu/eu.s3.tfbackend -backend-config=eu/dev/terraform.tfbackend")

	cmd = NewTerraformCommand("plan")
	err = cmd.Update(module, cwd)
	is.NoErr(err)
	is.Equal(cmd.String(), "TF_DATA_DIR=eu/dev/.terraform terraform -chdir=../.. plan -var-file=eu/eu.auto.tfvars -var-file=eu/dev/terraform.tfvars")
}
