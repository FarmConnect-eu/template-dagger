package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Destroy destroys infrastructure managed by Terraform
//
// WARNING: This will destroy all resources managed by the Terraform configuration.
// Make sure to backup any important data before running this command.
//
// This function runs `terraform init` followed by `terraform destroy`.
// Variables must be configured beforehand using WithVariable().
//
// Example usage:
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env://PM_API_URL --secret --tf-var \
//	  with-variable --key target_node --value pve-node-01 --tf-var \
//	  destroy --source=../infra-postgres --auto-approve
func (m *Terraform) Destroy(
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

	// Run terraform destroy
	args := []string{"terraform", "destroy"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	container = container.WithExec(args)

	// Return output
	return container.Stdout(ctx)
}
