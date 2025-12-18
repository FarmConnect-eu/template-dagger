package main

import "dagger/terraform/internal/dagger"

// WithVariable adds a non-secret variable to the module
func (m *Terraform) WithVariable(
	key string,
	value string,
	// +optional
	// +default=false
	tfVar bool,
) *Terraform {
	newVar := Variable{
		Key:         key,
		Value:       value,
		SecretValue: nil,
		TfVar:       tfVar,
	}

	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	newFiles := make([]*dagger.File, len(m.TfVarsFiles))
	copy(newFiles, m.TfVarsFiles)

	return &Terraform{
		Variables:        append(newVariables, newVar),
		State:            m.State,
		TerraformVersion: m.TerraformVersion,
		TfVarsFiles:      newFiles,
	}
}
