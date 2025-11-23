package main

import "dagger/terraform/internal/dagger"

// WithSecret adds un secret au module
func (m *Terraform) WithSecret(
	key string,
	value *dagger.Secret,
	// +optional
	// +default=false
	tfVar bool,
) *Terraform {
	newVar := Variable{
		Key:         key,
		Value:       "",
		SecretValue: value,
		TfVar:       tfVar,
	}

	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        append(newVariables, newVar),
		State:            m.State,
		TerraformVersion: m.TerraformVersion,
	}
}
