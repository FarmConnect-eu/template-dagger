package main

// WithTerraformVersion configure la version de Terraform à utiliser
//
// Par défaut, la version 1.9.8 est utilisée
//
// Exemple:
//
//	dagger call \
//	  with-terraform-version --version 1.10.0 \
//	  plan --source . --workdir terraform
func (m *Terraform) WithTerraformVersion(
	// Version de Terraform (ex: "1.9.8", "1.10.0", "latest")
	version string,
) *Terraform {
	// Deep copy pour éviter les mutations (pattern immutable)
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            m.State,
		TerraformVersion: version,
	}
}
