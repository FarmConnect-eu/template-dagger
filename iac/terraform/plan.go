package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Plan shows what changes will be made to infrastructure
//
// This function runs `terraform init` followed by `terraform plan`.
// Variables must be configured beforehand using WithVariable().
//
// Example usage:
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env://PM_API_URL --secret --tf-var \
//	  with-variable --key target_node --value pve-node-01 --tf-var \
//	  plan --source=../infra-postgres
func (m *Terraform) Plan(
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

	// Run terraform plan
	container = container.WithExec([]string{"terraform", "plan"})

	// Return output
	return container.Stdout(ctx)
}
