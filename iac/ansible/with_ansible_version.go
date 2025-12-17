package main

// WithAnsibleVersion sets the Ansible version (default: 11.1.0).
// Note: Ansible versioning changed at 3.0. Use 10.x/11.x for ansible-core 2.17+.
func (m *Ansible) WithAnsibleVersion(
	version string,
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
		AnsibleVersion: version,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		RolesPath:      m.RolesPath,
		Templates:      m.Templates,
		GroupVars:      m.GroupVars,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
