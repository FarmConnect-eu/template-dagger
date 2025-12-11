package main

// WithTarget sets the build target for multi-stage Dockerfiles
func (m *Docker) WithTarget(
	// Stage name from Dockerfile (e.g., "builder", "production")
	target string,
) *Docker {
	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs))
	copy(newBuildArgs, m.BuildArgs)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)

	return &Docker{
		BuildArgs:        newBuildArgs,
		Tags:             newTags,
		Target:           target,
		Platform:         m.Platform,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}
