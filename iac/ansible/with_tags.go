package main

// WithTags sets tags to run only specific tasks.
func (m *Ansible) WithTags(tags []string) *Ansible {
	newTags := make([]string, len(tags))
	copy(newTags, tags)

	return &Ansible{
		Variables:      m.Variables,
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		ExtraVars:      m.ExtraVars,
		Tags:           newTags,
		SkipTags:       m.SkipTags,
	}
}

// WithSkipTags sets tags to skip specific tasks.
func (m *Ansible) WithSkipTags(skipTags []string) *Ansible {
	newSkipTags := make([]string, len(skipTags))
	copy(newSkipTags, skipTags)

	return &Ansible{
		Variables:      m.Variables,
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		ExtraVars:      m.ExtraVars,
		Tags:           m.Tags,
		SkipTags:       newSkipTags,
	}
}
