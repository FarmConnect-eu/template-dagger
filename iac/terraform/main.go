// Package main provides a Dagger module for executing Terraform in containers.
package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

type Variable struct {
	Key         string
	Value       string
	SecretValue *dagger.Secret
	TfVar       bool
}

type StateConfig struct {
	Backend  string
	Bucket   string
	Key      string
	Region   string
	Endpoint string
}

type Terraform struct {
	Variables        []Variable
	State            *StateConfig
	TerraformVersion string
}

func New() *Terraform {
	return &Terraform{
		Variables:        []Variable{},
		State:            nil,
		TerraformVersion: "1.10.6",
	}
}

func (m *Terraform) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "Terraform module loaded successfully"}).
		Stdout(ctx)
}
