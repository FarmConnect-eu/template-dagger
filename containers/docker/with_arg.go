package main

// WithArg adds a Docker build argument (--build-arg)
//
// Build arguments are available during image build and can be referenced
// in Dockerfile using ARG instructions. Chain multiple calls for multiple args.
func (m *Docker) WithArg(
	// Argument name (e.g., "VERSION", "BUILD_DATE")
	key string,
	// Argument value
	value string,
) *Docker {
	newArg := DockerBuildArg{
		Key:   key,
		Value: value,
	}

	// Deep copy slice (immutable pattern)
	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs), len(m.BuildArgs)+1)
	copy(newBuildArgs, m.BuildArgs)

	// Deep copy tags slice
	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)

	return &Docker{
		BuildArgs:        append(newBuildArgs, newArg),
		Tags:             newTags,
		Target:           m.Target,
		Platform:         m.Platform,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}
