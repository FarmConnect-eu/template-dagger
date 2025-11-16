package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Deploy deploys infrastructure using Terraform
//
// This function runs `terraform init` followed by `terraform apply`.
// Variables must be configured beforehand using WithVariable().
//
// Example usage:
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env://PM_API_URL --secret --tf-var \
//	  with-variable --key target_node --value pve-node-01 --tf-var \
//	  deploy --source=../infra-postgres --auto-approve
func (m *Terraform) Deploy(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Working directory relative to source (default: terraform)
	// +optional
	// +default="terraform"
	workdir string,
	// Auto-approve without confirmation
	// +optional
	// +default=false
	autoApprove bool,
) (string, error) {
	if workdir == "" {
		workdir = "terraform"
	}

	// Configure backend if state is configured
	source, err := m.configureBackend(ctx, source, workdir)
	if err != nil {
		return "", err
	}

	// Build base container
	container := m.buildContainer(source, workdir)

	// Inject variables
	container, err = m.injectVariables(ctx, container)
	if err != nil {
		return "", err
	}

	// Run terraform init
	container = container.WithExec([]string{"terraform", "init"})

	// Run terraform apply
	args := []string{"terraform", "apply"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	container = container.WithExec(args)

	// Return output
	return container.Stdout(ctx)
}
