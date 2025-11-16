package main

// WithAnsibleVersion sets the Ansible version (default: "2.15")
// Installed on Ubuntu 22.04 from official PPA.
func (m *Ansible) WithAnsibleVersion(
	// +optional
	// +default="2.15"
	version string,
) *Ansible {
	if version == "" {
		version = "2.15"
	}

	// Deep copy to avoid mutation
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	return &Ansible{
		Variables:      newVariables,
		AnsibleVersion: version,
	}
}

