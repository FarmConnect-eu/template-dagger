package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
	"time"
)

// Deploy deploys the Docker Compose stack
//
// This function performs:
// 1. Pull and start containers with --pull always --force-recreate
// 2. Display container status
//
// The --pull always flag ensures images are always re-downloaded from registry,
// bypassing local cache. This guarantees "latest" tags get the actual latest version.
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

	// Deploy with force pull and recreate
	// --pull always: forces re-download of images (ignores local cache)
	// --force-recreate: recreates containers even if config unchanged
	//
	// IMPORTANT: WithEnvVariable with timestamp prevents Dagger from caching
	// the execution result. Without this, Dagger may return cached output
	// without actually running the deployment commands on the remote host.
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
