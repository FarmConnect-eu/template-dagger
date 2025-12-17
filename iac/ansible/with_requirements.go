package main

import "dagger/ansible/internal/dagger"

// WithRequirements sets the Ansible requirements file (requirements.yml for collections/roles).
func (m *Ansible) WithRequirements(
	requirements *dagger.File,
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
		Inventory:      m.Inventory,
		Requirements:   requirements,
		RolesPath:      m.RolesPath,
		Templates:      m.Templates,
		GroupVars:      m.GroupVars,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
