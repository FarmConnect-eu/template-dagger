package main

// WithTag adds image tags.
//
// If the tag is a semantic version (e.g., "v1.2.3", "v1.0.0-rc4"), it automatically
// generates all appropriate tags:
//   - v1.0.0-rc4 -> v1.0.0-rc4, v1.0.0, v1.0, v1, rc
//   - v1.0.0-release -> v1.0.0-release, v1.0.0, v1.0, v1, release, latest
//   - v1.0.0 -> v1.0.0, v1.0, v1, latest
//
// For non-semver tags (e.g., "dev", "latest"), only that tag is added.
func (m *Docker) WithTag(
	// Image tag (e.g., "v1.2.3", "v1.0.0-rc4", "latest", "dev")
	tag string,
) *Docker {
	// Parse semantic version tags
	tags := parseSemanticVersion(tag)

	// Deep copy build args slice
	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs))
	copy(newBuildArgs, m.BuildArgs)

	// Deep copy tags slice (immutable pattern)
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
