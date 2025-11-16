package main

// WithVariable adds a variable. Supports literal, env://, file:// values.
func (m *Ansible) WithVariable(
	key string,
	value string,
	// +optional
	// +default=false
	secret bool,
) *Ansible {
	newVar := Variable{
		Key:      key,
		Value:    value,
		IsSecret: secret,
	}

	// Deep copy to avoid mutation
	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	return &Ansible{
		Variables:      append(newVariables, newVar),
		AnsibleVersion: m.AnsibleVersion,
	}
}
