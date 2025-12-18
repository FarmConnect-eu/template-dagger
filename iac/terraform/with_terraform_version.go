package main

import "dagger/terraform/internal/dagger"

// WithTerraformVersion configures the Terraform/OpenTofu version to use (default: 1.10.6)
func (m *Terraform) WithTerraformVersion(
	// Terraform version (e.g., "1.9.8", "1.10.0", "latest")
	version string,
) *Terraform {
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	newFiles := make([]*dagger.File, len(m.TfVarsFiles))
	copy(newFiles, m.TfVarsFiles)

	return &Terraform{
		Variables:        newVariables,
		State:            m.State,
		TerraformVersion: version,
		TfVarsFiles:      newFiles,
	}
}
