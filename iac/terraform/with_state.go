package main

import "dagger/terraform/internal/dagger"

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

	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	newFiles := make([]*dagger.File, len(m.TfVarsFiles))
	copy(newFiles, m.TfVarsFiles)

	return &Terraform{
		Variables:        newVariables,
		State:            stateConfig,
		TerraformVersion: m.TerraformVersion,
		TfVarsFiles:      newFiles,
	}
}
