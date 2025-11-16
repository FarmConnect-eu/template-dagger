// Terraform CI/CD module for Dagger.
// Automates Terraform workflows with secret management and flexible source input.

package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Variable to inject into Terraform (supports literal, env://, file://)
type Variable struct {
	Key      string
	Value    string
	IsSecret bool
	TfVar    bool
}

// StateConfig for Terraform backend
type StateConfig struct {
	Backend string
	Bucket  string
	Key     string
	Region  string
}

// Terraform module
type Terraform struct {
	Variables        []Variable
	State            *StateConfig
	TerraformVersion string
}

// Test module health
func (m *Terraform) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "Terraform CI module loaded successfully"}).
		Stdout(ctx)
}
