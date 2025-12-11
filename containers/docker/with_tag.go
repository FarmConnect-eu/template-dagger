package main

// WithTag adds image tags.
// For semantic versions (v1.2.3), automatically generates hierarchy tags (v1.2, v1, latest).
func (m *Docker) WithTag(
	// Image tag (e.g., "v1.2.3", "v1.0.0-rc4", "latest", "dev")
	tag string,
) *Docker {
	tags := parseSemanticVersion(tag)

	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs))
	copy(newBuildArgs, m.BuildArgs)

	newTags := make([]string, len(m.Tags), len(m.Tags)+len(tags))
	copy(newTags, m.Tags)
	newTags = append(newTags, tags...)

	return &Docker{
		BuildArgs:        newBuildArgs,
		Tags:             newTags,
		Target:           m.Target,
		Platform:         m.Platform,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}
