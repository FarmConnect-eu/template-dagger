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
	c := m.clone()
	c.RegistryHost = host
	c.RegistryUsername = username
	c.RegistryPassword = password
	return c
}
