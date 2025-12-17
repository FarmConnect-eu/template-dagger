package main

import (
	"fmt"

	"dagger/docker/internal/dagger"
)

// Build builds a Docker image from a Dockerfile using Dagger's native build
//
// Returns a container that can be exported or pushed to a registry.
func (m *Docker) Build(
	// Directory containing Dockerfile and build context
	source *dagger.Directory,
	// Image name for reference (e.g., "myapp")
	imageName string,
	// Path to Dockerfile relative to source
	// +optional
	// +default="Dockerfile"
	dockerfile string,
) (*dagger.Container, error) {
	if err := validateImageName(imageName); err != nil {
		return nil, fmt.Errorf("invalid image name: %w", err)
	}

	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}

	// Build options
	buildOpts := dagger.DirectoryDockerBuildOpts{
		Dockerfile: dockerfile,
		Platform:   m.Platform,
	}

	// Set target if configured
	if m.Target != "" {
		buildOpts.Target = m.Target
	}

	// Add build arguments (convert to Dagger BuildArg type)
	if len(m.BuildArgs) > 0 {
		buildArgs := make([]dagger.BuildArg, 0, len(m.BuildArgs))
		for _, arg := range m.BuildArgs {
			buildArgs = append(buildArgs, dagger.BuildArg{
				Name:  arg.Key,
				Value: arg.Value,
			})
		}
		buildOpts.BuildArgs = buildArgs
	}

	// Build container using Directory.DockerBuild()
	container := source.DockerBuild(buildOpts)

	// Add image reference as label
	tags := m.getDefaultTags()
	for _, tag := range tags {
		fullRef := buildFullReference(m.RegistryHost, imageName, tag)
		container = container.WithLabel("org.opencontainers.image.ref.name", fullRef)
	}

	return container, nil
}
