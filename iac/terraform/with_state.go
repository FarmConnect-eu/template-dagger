package main

// WithState configures Terraform backend (s3, gcs, azurerm, local)
func (m *Terraform) WithState(
	backend string,
	// +optional
	bucket string,
	// +optional
	key string,
	// +optional
	region string,
) *Terraform {
	stateConfig := &StateConfig{
		Backend: backend,
		Bucket:  bucket,
		Key:     key,
		Region:  region,
	}

	// Deep copy to avoid mutation
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            stateConfig,
		TerraformVersion: m.TerraformVersion,
	}
}
