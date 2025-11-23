package main

// WithTerraformVersion configures la version de Terraform à utiliser
//
// Par défaut, la version 1.9.8 est utilisée
//
func (m *Terraform) WithTerraformVersion(
	// Version de Terraform (ex: "1.9.8", "1.10.0", "latest")
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
