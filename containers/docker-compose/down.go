package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Down stops and removes Docker Compose containers
//
// This function stops all containers and removes them, networks, and volumes
// created by the Docker Compose stack.
//
// Parameters:
//   - source: Directory containing the docker-compose.yml file
//   - composePath: Path to docker-compose.yml relative to source (default: "docker-compose.yml")
//   - projectName: Docker Compose project name (optional, uses directory name if not set)
//
// Example:
//
//	dagger call down \
//	  --source . \
//	  --compose-path docker/docker-compose.yml \
//	  --project-name chat
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

	// Add project name if specified
	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	// Stop and remove containers
	downCmd := append(composeCmd, "down")
	output, err := container.WithExec(downCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to stop containers: %w", err)
	}

	return fmt.Sprintf("Containers stopped successfully\n\n%s", output), nil
}
