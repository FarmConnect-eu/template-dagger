package main

import "dagger/ansible/internal/dagger"

// WithInventory sets the inventory file.
func (m *Ansible) WithInventory(
	inventory *dagger.File,
) *Ansible {
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	newExtraVars := make([]KeyValue, len(m.ExtraVars))
	copy(newExtraVars, m.ExtraVars)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)
	newSkipTags := make([]string, len(m.SkipTags))
	copy(newSkipTags, m.SkipTags)

	return &Ansible{
		Variables:      newVariables,
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      inventory,
		Requirements:   m.Requirements,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
