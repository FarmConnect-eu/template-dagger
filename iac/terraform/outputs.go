package main

import (
	"context"
	"dagger/terraform/internal/dagger"
)

// Outputs retrieves Terraform outputs from infrastructure deployment
//
// Returns all outputs as JSON, or a specific module/output if specified.
// Variables must be configured beforehand using WithVariable().
//
// Example: Get all outputs as JSON
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env://PM_API_URL --secret --tf-var \
//	  outputs --source=../infra-nfs
//
// Example: Get specific output
//
//	dagger call \
//	  with-variable --key proxmox_api_url --value env://PM_API_URL --secret --tf-var \
//	  outputs --source=../infra-postgres --output-name postgres_ip
func (m *Terraform) Outputs(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Working directory relative to source (default: terraform)
	// +optional
	// +default="terraform"
	workdir string,
	// Specific output name (returns all as JSON if empty)
	// +optional
	outputName string,
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

	// Run terraform output
	args := []string{"terraform", "output"}
	if outputName != "" {
		args = append(args, outputName)
	} else {
		args = append(args, "-json")
	}
	container = container.WithExec(args)

	// Return output
	return container.Stdout(ctx)
}
