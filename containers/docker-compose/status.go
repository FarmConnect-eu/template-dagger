package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Status displays the status of Docker Compose containers
func (m *DockerCompose) Status(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="docker-compose.yml"
	composePath string,
	// +optional
	projectName string,
) (string, error) {
	if composePath == "" {
		composePath = "docker-compose.yml"
	}

	container := m.buildContainer(ctx, source, composePath)
	composeCmd := getComposeCommand(composePath)

	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	psCmd := append(composeCmd, "ps")
	output, err := container.WithExec(psCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	return output, nil
}
