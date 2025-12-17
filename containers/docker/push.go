package main

import (
	"context"
	"fmt"

	"dagger/docker/internal/dagger"
)

// Push pushes a built container to a registry
//
// Requires registry authentication configured via WithRegistry().
// Pushes all configured tags (defaults to "latest" if none specified).
// Returns the image digest.
func (m *Docker) Push(
	ctx context.Context,
	// Built container from Build()
	container *dagger.Container,
	// Image name without registry prefix (e.g., "myapp" or "myorg/myapp")
	imageName string,
) (string, error) {
	if m.RegistryHost == "" {
		return "", fmt.Errorf("registry not configured: use WithRegistry() first")
	}
	if m.RegistryUsername == "" || m.RegistryPassword == nil {
		return "", fmt.Errorf("registry credentials missing: use WithRegistry() with username and password")
	}
	if err := validateImageName(imageName); err != nil {
		return "", fmt.Errorf("invalid image name: %w", err)
	}

	// Authenticate to registry
	container = container.WithRegistryAuth(
		m.RegistryHost,
		m.RegistryUsername,
		m.RegistryPassword,
	)

	// Push all configured tags
	tags := m.getDefaultTags()
	var lastDigest string

	for _, tag := range tags {
		fullReference := buildFullReference(m.RegistryHost, imageName, tag)
		digest, err := container.Publish(ctx, fullReference)
		if err != nil {
			return "", fmt.Errorf("failed to push %s: %w", fullReference, err)
		}
		lastDigest = digest
	}

	return lastDigest, nil
}
