package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithRegistry configures Docker registry authentication
//
// # This allows pulling from private registries before deploying
//
// Parameters:
//   - host: Registry hostname (e.g., "registry.example.com", "ghcr.io")
//   - username: Registry username
//   - password: Registry password
//
// Example:
//
//	dagger call with-registry \
//	  --host ghcr.io \
//	  --username env:GITHUB_USER \
//	  --password env:GITHUB_TOKEN
func (m *DockerCompose) WithRegistry(
	host string,
	username string,
	password *dagger.Secret,
) *DockerCompose {
	return &DockerCompose{
		RegistryHost:     host,
		RegistryUsername: username,
		RegistryPassword: password,
		Variables:        copyVariables(m.Variables),
		SSHHost:          m.SSHHost,
		SSHUser:          m.SSHUser,
		SSHPort:          m.SSHPort,
		SSHKey:           m.SSHKey,
		EnvFile:          m.EnvFile,
	}
}
