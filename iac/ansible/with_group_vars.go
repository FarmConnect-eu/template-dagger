package main

import "dagger/ansible/internal/dagger"

// WithGroupVars sets the group_vars directory to be mounted in the container.
// Group vars will be available at /work/group_vars for Ansible to use.
func (m *Ansible) WithGroupVars(
	// Directory containing Ansible group_vars
	groupVars *dagger.Directory,
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
		Requirements:   m.Requirements,
		RolesPath:      m.RolesPath,
		Templates:      m.Templates,
		GroupVars:      groupVars,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
