package main

import "dagger/ansible/internal/dagger"

// WithTemplate sets the templates directory to be mounted in the container.
// Templates will be available at /work/templates for Ansible to use.
func (m *Ansible) WithTemplate(
	// Directory containing Ansible templates
	templates *dagger.Directory,
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
		Templates:      templates,
		GroupVars:      m.GroupVars,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
