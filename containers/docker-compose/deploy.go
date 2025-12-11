package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Deploy deploys the Docker Compose stack (pull, recreate, status)
func (m *DockerCompose) Deploy(
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

	pullCmd := append(composeCmd, "pull")
	container = container.WithExec(pullCmd)

	upCmd := append(composeCmd, "up", "-d", "--force-recreate")
	container = container.WithExec(upCmd)

	psCmd := append(composeCmd, "ps")
	output, err := container.WithExec(psCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("deployment failed: %w", err)
	}

	return fmt.Sprintf("Deployment successful\n\n%s", output), nil
}
