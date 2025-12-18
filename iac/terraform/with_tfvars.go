package main

import "dagger/terraform/internal/dagger"

// WithTfVarsFile adds a tfvars file to be used during terraform operations
func (m *Terraform) WithTfVarsFile(
	// The tfvars file to mount
	file *dagger.File,
) *Terraform {
	newFiles := make([]*dagger.File, len(m.TfVarsFiles), len(m.TfVarsFiles)+1)
	copy(newFiles, m.TfVarsFiles)

	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            m.State,
		TerraformVersion: m.TerraformVersion,
		TfVarsFiles:      append(newFiles, file),
	}
}
