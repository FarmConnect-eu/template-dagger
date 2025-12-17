package main

import (
	"context"
	"fmt"

	"dagger/docker-compose/internal/dagger"
)

// buildContainer creates a container with Docker Compose installed and configured
func (m *DockerCompose) buildContainer(
	ctx context.Context,
	source *dagger.Directory,
	composePath string,
) *dagger.Container {
	// Start with Docker CLI image
	container := dag.Container().
		From("docker:27-cli").
		WithMountedDirectory("/workspace", source).
		WithWorkdir("/workspace")

	// Configure SSH context if provided, otherwise use local socket
	if m.SSHHost != "" && m.SSHKey != nil {
		// Install Docker Compose and SSH client for remote deployment
		container = container.WithExec([]string{
			"sh", "-c",
			"apk add --no-cache docker-cli-compose openssh-client",
		})

		// Setup SSH directory and key
		// Mount secret to temp location, then copy with correct permissions (mounted secrets are read-only)
		container = container.
			WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
			WithMountedSecret("/tmp/ssh_key", m.SSHKey).
			WithExec([]string{"sh", "-c", "cp /tmp/ssh_key /root/.ssh/id_ed25519 && chmod 600 /root/.ssh/id_ed25519"})

		// Configure SSH to skip host key verification (for CI/CD)
		sshConfig := "Host *\n  StrictHostKeyChecking no\n  UserKnownHostsFile /dev/null\n  LogLevel ERROR\n"
		container = container.WithNewFile("/root/.ssh/config", sshConfig, dagger.ContainerWithNewFileOpts{
			Permissions: 0o600,
		})

		// Set DOCKER_HOST to SSH endpoint
		sshHost := fmt.Sprintf("ssh://%s@%s", m.SSHUser, m.SSHHost)
		if m.SSHPort != 22 {
			sshHost = fmt.Sprintf("ssh://%s@%s:%d", m.SSHUser, m.SSHHost, m.SSHPort)
		}
		container = container.WithEnvVariable("DOCKER_HOST", sshHost)
	} else {
		// Local Docker socket not supported in current Dagger SDK
		// SSH context is required for docker-compose deployments
		// Use WithContext to configure SSH connection
		container = container.WithExec([]string{
			"sh", "-c",
			"apk add --no-cache docker-cli-compose",
		})
	}

	// Mount .env file if provided
	if m.EnvFile != nil {
		container = container.WithMountedFile("/workspace/.env", m.EnvFile)
	}

	// Authenticate to registry on remote Docker host via SSH
	if m.RegistryHost != "" && m.RegistryUsername != "" && m.RegistryPassword != nil {
		container = container.
			WithEnvVariable("REGISTRY_HOST", m.RegistryHost).
			WithEnvVariable("REGISTRY_USERNAME", m.RegistryUsername).
			WithSecretVariable("REGISTRY_PASSWORD", m.RegistryPassword).
			WithExec([]string{
				"sh", "-c",
				"echo $REGISTRY_PASSWORD | docker login $REGISTRY_HOST -u $REGISTRY_USERNAME --password-stdin",
			})
	}

	// Inject environment variables
	for _, v := range m.Variables {
		if v.Secret != nil {
			container = container.WithSecretVariable(v.Key, v.Secret)
		} else {
			container = container.WithEnvVariable(v.Key, v.Value)
		}
	}

	return container
}

// getComposeCommand returns the docker compose command with the compose file path
func getComposeCommand(composePath string) []string {
	if composePath == "" {
		composePath = "docker-compose.yml"
	}
	return []string{"docker", "compose", "-f", composePath}
}
