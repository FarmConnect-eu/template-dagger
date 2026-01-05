package main

import (
	"context"
	"fmt"
	"time"

	"dagger/docker-compose/internal/dagger"
)

// Deploy deploys the Docker Compose stack
//
// This function performs:
// 1. Pull latest images
// 2. Stop and remove existing containers
// 3. Start new containers with --force-recreate
// 4. Display container status
//
// Parameters:
//   - source: Directory containing the docker-compose.yml file
//   - composePath: Path to docker-compose.yml relative to source (default: "docker-compose.yml")
//   - projectName: Docker Compose project name (optional, uses directory name if not set)
//
// Example:
//
//	dagger call deploy \
//	  --source . \
//	  --compose-path docker/docker-compose.yml \
//	  --project-name chat
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

	// Add project name if specified
	if projectName != "" {
		composeCmd = append(composeCmd, "-p", projectName)
	}

	// Pull latest images
	pullCmd := append(composeCmd, "pull")
	container = container.WithExec(pullCmd)

	// Deploy with force recreate
	upCmd := append(composeCmd, "up", "-d", "--pull", "always", "--force-recreate")
	container = container.
		WithEnvVariable("DAGGER_CACHE_BUSTER", time.Now().String()).
		WithExec(upCmd)

	// Get container status
	psCmd := append(composeCmd, "ps")
	output, err := container.WithExec(psCmd).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("deployment failed: %w", err)
	}

	return fmt.Sprintf("Deployment successful\n\n%s", output), nil
}
