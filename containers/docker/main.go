// Docker module for building and pushing container images
//
// This module provides Docker image build and push operations using Dagger's
// native container build capabilities. Supports Azure ACR and other registries.
package main

import (
	"dagger/docker/internal/dagger"
)

// Docker module for building and pushing container images to registries
type Docker struct {
	// Registry authentication
	RegistryHost     string
	RegistryUsername string
	RegistryPassword *dagger.Secret

	// Build configuration
	BuildArgs []DockerBuildArg
	Tags      []string
	Target    string
	Platform  dagger.Platform
}

// DockerBuildArg represents a Docker build argument
type DockerBuildArg struct {
	Key   string
	Value string
}

// New creates a new Docker instance with default configuration
func New() *Docker {
	return &Docker{
		BuildArgs: []DockerBuildArg{},
		Tags:      []string{},
		Platform:  "linux/amd64",
	}
}

// Test verifies the module loads correctly
func (m *Docker) Test() string {
	return "Docker module loaded successfully"
}
