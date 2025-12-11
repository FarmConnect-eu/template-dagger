package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Logs retrieves logs from Docker Compose containers
func (m *DockerCompose) Logs(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="docker-compose.yml"
	composePath string,
	// +optional
	// +default=100
	tail int,
	// +optional
	projectName string,
) (string, error) {
	if composePath == "" {
		composePath = "docker-compose.yml"
	}
	if tail == 0 {
		tail = 100
	}

	container := m.buildContainer(ctx, source, composePath)
	composeCmd := getComposeCommand(composePath)

	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	logsCmd := append(composeCmd, "logs", "--tail", fmt.Sprintf("%d", tail))
	output, err := container.WithExec(logsCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve logs: %w", err)
	}

	return output, nil
}
