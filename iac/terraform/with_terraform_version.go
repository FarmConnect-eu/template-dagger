package main

// WithTerraformVersion sets the Terraform CLI version (default: "1.9")
func (m *Terraform) WithTerraformVersion(
	// +optional
	// +default="1.9"
	version string,
) *Terraform {
	if version == "" {
		version = "1.9"
	}

	return &Terraform{
		Variables:        m.Variables,
		State:            m.State,
		TerraformVersion: version,
	}
}
