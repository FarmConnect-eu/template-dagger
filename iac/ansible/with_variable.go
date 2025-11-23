package main

// WithVariable adds a non-secret variable (use WithSecret for secrets).
func (m *Ansible) WithVariable(
	key string,
	value string,
) *Ansible {
	newVar := Variable{
		Key:         key,
		Value:       value,
		SecretValue: nil,
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

// WithExtraVar adds an extra variable for Ansible (--extra-vars).
func (m *Ansible) WithExtraVar(
	key string,
	value string,
) *Ansible {
	newVariables := make([]Variable, len(m.Variables))
	copy(newVariables, m.Variables)

	newExtraVars := make([]KeyValue, len(m.ExtraVars), len(m.ExtraVars)+1)
	copy(newExtraVars, m.ExtraVars)
	newExtraVars = append(newExtraVars, KeyValue{Key: key, Value: value})

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)
	newSkipTags := make([]string, len(m.SkipTags))
	copy(newSkipTags, m.SkipTags)

	return &Ansible{
		Variables:      newVariables,
		AnsibleVersion: m.AnsibleVersion,
		Inventory:      m.Inventory,
		Requirements:   m.Requirements,
		ExtraVars:      newExtraVars,
		Tags:           newTags,
		SkipTags:       newSkipTags,
	}
}
