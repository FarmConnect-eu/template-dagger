package main

// WithState configures Terraform backend for state management (supports s3, gcs, azurerm, local).
func (m *Terraform) WithState(
	// Backend type (s3, gcs, azurerm, local)
	backend string,
	// Bucket/container name (unused for local)
	// +optional
	bucket string,
	// State key/file path
	// +optional
	key string,
	// Region (used for S3)
	// +optional
	region string,
	// S3-compatible endpoint URL (for MinIO, etc.)
	// +optional
	endpoint string,
) *Terraform {
	stateConfig := &StateConfig{
		Backend:  backend,
		Bucket:   bucket,
		Key:      key,
		Region:   region,
		Endpoint: endpoint,
	}

	// Deep copy pour éviter les mutations (pattern immutable)
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Terraform{
		Variables:        newVariables,
		State:            stateConfig,
		TerraformVersion: m.TerraformVersion,
	}
}
