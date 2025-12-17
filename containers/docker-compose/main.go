// Docker Compose module for Dagger
//
// This module provides Docker Compose operations with chainable configuration
// for registry authentication, environment variables, and deployment management.
//
// Example usage:
//
//	dagger call -m containers/docker-compose \
//	  with-context --host 172.16.24.97 --user admin --ssh-key env:SSH_KEY \
//	  with-registry --host registry.example.com --username env:USER --password env:PASS \
//	  with-secret --key DB_PASSWORD --value env:DB_PASSWORD \
//	  with-variable --key IMAGE_TAG --value v1.0.0 \
//	  deploy --source . --compose-path docker/docker-compose.yml --project-name myapp

package main

import (
	"dagger/docker-compose/internal/dagger"
)

// DockerCompose module for managing Docker Compose deployments
type DockerCompose struct {
	// Registry authentication configuration
	RegistryHost     string
	RegistryUsername string
	RegistryPassword *dagger.Secret

	// Environment variables to inject
	Variables []*Variable

	// SSH Context configuration for remote deployment
	SSHHost string
	SSHUser string
	SSHPort int
	SSHKey  *dagger.Secret

	// Environment file
	EnvFile *dagger.File
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
		Variables: []*Variable{},
		SSHPort:   22,
	}
}
