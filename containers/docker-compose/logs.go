package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

// Logs retrieves logs from Docker Compose containers
//
// This function fetches the logs from all containers in the Docker Compose stack.
//
// Parameters:
//   - source: Directory containing the docker-compose.yml file
//   - composePath: Path to docker-compose.yml relative to source (default: "docker-compose.yml")
//   - tail: Number of lines to show from the end of logs (default: 100)
//   - projectName: Docker Compose project name (optional, uses directory name if not set)
//
// Example:
//
//	dagger call logs \
//	  --source . \
//	  --compose-path docker/docker-compose.yml \
//	  --tail 50 \
//	  --project-name chat
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

	// Add project name if specified
	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	// Get logs
	logsCmd := append(composeCmd, "logs", "--tail", fmt.Sprintf("%d", tail))
	output, err := container.WithExec(logsCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve logs: %w", err)
	}

	return output, nil
}
