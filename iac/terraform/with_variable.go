package main

// WithVariable adds a variable. Supports literal, env://, file:// values.
func (m *Terraform) WithVariable(
	key string,
	value string,
	// +optional
	// +default=false
	secret bool,
	// +optional
	// +default=false
	tfVar bool,
) *Terraform {
	newVar := Variable{
		Key:      key,
		Value:    value,
		IsSecret: secret,
		TfVar:    tfVar,
	}

	// Deep copy to avoid mutation
	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        append(newVariables, newVar),
		State:            m.State,
		TerraformVersion: m.TerraformVersion,
	}
}
