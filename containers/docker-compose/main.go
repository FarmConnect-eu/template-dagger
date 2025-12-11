// Docker Compose module for Dagger
package main

import (
	"dagger/docker-compose/internal/dagger"
)

// DockerCompose module for managing Docker Compose deployments
type DockerCompose struct {
	// Registry authentication configuration
	RegistryHost     string
	RegistryUsername *dagger.Secret
	RegistryPassword *dagger.Secret

	// Environment variables to inject
	Variables map[string]*Variable
}

// Variable represents an environment variable with optional secret flag
type Variable struct {
	Key    string
	Value  string
	Secret *dagger.Secret
}

// New creates a new DockerCompose instance
func New() *DockerCompose {
	return &DockerCompose{
		Variables: make(map[string]*Variable),
	}
}
