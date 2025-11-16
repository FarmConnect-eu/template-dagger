package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Format formats the Terraform configuration files
//
// Runs terraform fmt to ensure consistent formatting.
// No variables needed for formatting.
//
// Example usage:
//
//	dagger call format --source=../infra-postgres
//	dagger call format --source=../infra-nfs --check
func (m *Terraform) Format(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Working directory relative to source (default: terraform)
	// +optional
	// +default="terraform"
	workdir string,
	// Check if files are formatted (returns error if not)
	// +optional
	// +default=false
	check bool,
) (*dagger.Directory, error) {
	if workdir == "" {
		workdir = "terraform"
	}

	// Build base container (no variables needed for formatting)
	container := m.buildContainer(source, workdir)

	// Run terraform fmt
	args := []string{"terraform", "fmt"}
	if check {
		args = append(args, "-check")
	}
	container = container.WithExec(args)

	// Return formatted directory
	return container.Directory("/work"), nil
}
