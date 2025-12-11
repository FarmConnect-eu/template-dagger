package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithRegistry configures Docker registry authentication
func (m *DockerCompose) WithRegistry(
	host string,
	username *dagger.Secret,
	password *dagger.Secret,
) *DockerCompose {
	m.RegistryHost = host
	m.RegistryUsername = username
	m.RegistryPassword = password
	return m
}
