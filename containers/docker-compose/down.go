package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Down stops and removes Docker Compose containers
func (m *DockerCompose) Down(
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

	downCmd := append(composeCmd, "down")
	output, err := container.WithExec(downCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to stop containers: %w", err)
	}

	return fmt.Sprintf("Containers stopped successfully\n\n%s", output), nil
}
