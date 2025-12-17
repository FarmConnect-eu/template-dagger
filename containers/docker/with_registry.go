package main

import (
	"dagger/docker/internal/dagger"
)

// WithRegistry configures Docker registry authentication for pushing images
//
// Supports Azure ACR (*.azurecr.io) and other Docker-compatible registries.
// Use environment variable references (env:VAR_NAME) for credentials.
func (m *Docker) WithRegistry(
	// Registry hostname (e.g., "myregistry.azurecr.io")
	host string,
	// Registry username (use env:VAR_NAME for environment variables)
	username string,
	// Registry password or token (use env:VAR_NAME for environment variables)
	password *dagger.Secret,
) *Docker {
	// Deep copy slices
	newBuildArgs := make([]DockerBuildArg, len(m.BuildArgs))
	copy(newBuildArgs, m.BuildArgs)

	newTags := make([]string, len(m.Tags))
	copy(newTags, m.Tags)

	return &Docker{
		RegistryHost:     host,
		RegistryUsername: username,
		RegistryPassword: password,
		BuildArgs:        newBuildArgs,
		Tags:             newTags,
		Target:           m.Target,
		Platform:         m.Platform,
	}
}
