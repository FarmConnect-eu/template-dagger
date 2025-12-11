package main

import "dagger/ansible/internal/dagger"

// WithRolesPath sets an external roles directory to be mounted in the container.
// The roles will be available at /work/roles and ANSIBLE_ROLES_PATH will be set.
func (m *Ansible) WithRolesPath(
	// Directory containing Ansible roles
	path *dagger.Directory,
) *Ansible {
	newM := &Ansible{
		Variables:      make([]Variable, len(m.Variables)),
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		RolesPath:      path,
		ExtraVars:      make([]KeyValue, len(m.ExtraVars)),
		Tags:           make([]string, len(m.Tags)),
		SkipTags:       make([]string, len(m.SkipTags)),
	}

	copy(newM.Variables, m.Variables)
	copy(newM.ExtraVars, m.ExtraVars)
	copy(newM.Tags, m.Tags)
	copy(newM.SkipTags, m.SkipTags)

	return newM
}
