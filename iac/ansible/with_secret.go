package main

import "dagger/ansible/internal/dagger"

// WithSecret adds a secret variable (supports env:, file: prefixes).
func (m *Ansible) WithSecret(
	key string,
	value *dagger.Secret,
) *Ansible {
	newVar := Variable{
		Key:         key,
		Value:       "",
		SecretValue: value,
	}

	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+1)
	copy(newVariables, m.Variables)

	newExtraVars := make([]KeyValue, len(m.ExtraVars))
	copy(newExtraVars, m.ExtraVars)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)
	newSkipTags := make([]string, len(m.SkipTags))
	copy(newSkipTags, m.SkipTags)

	return &Ansible{
		Variables:      append(newVariables, newVar),
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
