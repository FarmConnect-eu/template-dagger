package main

// WithTerraformVersion configure la version de Terraform
func (m *Terraform) WithTerraformVersion(
	// +optional
	// +default="1.9.8"
	version string,
) *Terraform {
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            m.State,
		TerraformVersion: version,
	}
}
