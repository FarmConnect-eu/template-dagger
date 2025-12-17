package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Status displays the status of Docker Compose containers
//
// This function shows the current state of all containers in the Docker Compose stack,
// including their names, status, ports, and health status.
//
// Parameters:
//   - source: Directory containing the docker-compose.yml file
//   - composePath: Path to docker-compose.yml relative to source (default: "docker-compose.yml")
//   - projectName: Docker Compose project name (optional, uses directory name if not set)
//
// Example:
//
//	dagger call status \
//	  --source . \
//	  --compose-path docker/docker-compose.yml \
//	  --project-name chat
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

	// Add project name if specified
	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	// Get container status
	psCmd := append(composeCmd, "ps")
	output, err := container.WithExec(psCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	return output, nil
}
