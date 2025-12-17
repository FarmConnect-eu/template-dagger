package main

import (
	"context"
	"fmt"

	"dagger/docker/internal/dagger"
)

// BuildAndPush builds and pushes a Docker image in one operation
//
// Convenience function combining Build() and Push().
// Requires registry authentication configured via WithRegistry().
// Returns the image digest.
func (m *Docker) BuildAndPush(
	ctx context.Context,
	// Directory containing Dockerfile and build context
	source *dagger.Directory,
	// Image name without registry prefix (e.g., "myapp")
	imageName string,
	// Path to Dockerfile relative to source
	// +optional
	// +default="Dockerfile"
	dockerfile string,
) (string, error) {
	// Validate registry early (fail fast)
	if m.RegistryHost == "" {
		return "", fmt.Errorf("registry not configured: use WithRegistry() before BuildAndPush()")
	}

	// Build
	container, err := m.Build(source, imageName, dockerfile)
	if err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	// Push
	digest, err := m.Push(ctx, container, imageName)
	if err != nil {
		return "", fmt.Errorf("push failed: %w", err)
	}

	return digest, nil
}
