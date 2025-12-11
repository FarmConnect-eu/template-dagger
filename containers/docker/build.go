package main

import (
	"fmt"

	"dagger/docker/internal/dagger"
)

// Build builds a Docker image from a Dockerfile using Dagger's native build
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

	buildOpts := dagger.DirectoryDockerBuildOpts{
		Dockerfile: dockerfile,
		Platform:   m.Platform,
	}

	if m.Target != "" {
		buildOpts.Target = m.Target
	}

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

	container := source.DockerBuild(buildOpts)

	tags := m.getDefaultTags()
	for _, tag := range tags {
		fullRef := buildFullReference(m.RegistryHost, imageName, tag)
		container = container.WithLabel("org.opencontainers.image.ref.name", fullRef)
	}

	return container, nil
}
