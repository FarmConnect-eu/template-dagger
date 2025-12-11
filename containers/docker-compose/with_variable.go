package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithVariable adds an environment variable to inject into Docker Compose
func (m *DockerCompose) WithVariable(
	key string,
	value string,
	// +optional
	secret *dagger.Secret,
) *DockerCompose {
	m.Variables[key] = &Variable{
		Key:    key,
		Value:  value,
		Secret: secret,
	}
	return m
}
