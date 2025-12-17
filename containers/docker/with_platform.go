package main

import (
	"dagger/docker/internal/dagger"
)

// WithPlatform sets the target platform for the build
//
// Overrides the default platform (linux/amd64).
func (m *Docker) WithPlatform(
	// Target platform (e.g., "linux/amd64", "linux/arm64")
	platform dagger.Platform,
) *Docker {
	// Deep copy slices
	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs))
	copy(newBuildArgs, m.BuildArgs)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)

	return &Docker{
		BuildArgs:        newBuildArgs,
		Tags:             newTags,
		Target:           m.Target,
		Platform:         platform,
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
	}
}
