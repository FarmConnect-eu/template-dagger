package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Validate validates the Terraform configuration
//
// Checks for syntax errors and configuration issues without connecting to Proxmox.
// No variables needed for validation.
//
// Example usage:
//
//	dagger call validate --source=../infra-traefik
func (m *Terraform) Validate(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Working directory relative to source (default: terraform)
	// +optional
	// +default="terraform"
	workdir string,
) (string, error) {
	if workdir == "" {
		workdir = "terraform"
	}

	// Build base container (no variables needed for validation)
	container := m.buildContainer(source, workdir)

	// Run terraform init
	container = container.WithExec([]string{"terraform", "init", "-backend=false"})

	// Run terraform validate
	container = container.WithExec([]string{"terraform", "validate"})

	// Return output
	return container.Stdout(ctx)
}
