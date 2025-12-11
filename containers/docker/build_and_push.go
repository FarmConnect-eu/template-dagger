package main

import (
	"context"
	"fmt"

	"dagger/docker/internal/dagger"
)

// BuildAndPush builds and pushes a Docker image in one operation
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
	if m.RegistryHost == "" {
		return "", fmt.Errorf("registry not configured: use WithRegistry() before BuildAndPush()")
	}

	container, err := m.Build(source, imageName, dockerfile)
	if err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	digest, err := m.Push(ctx, container, imageName)
	if err != nil {
		return "", fmt.Errorf("push failed: %w", err)
	}

	return digest, nil
}
