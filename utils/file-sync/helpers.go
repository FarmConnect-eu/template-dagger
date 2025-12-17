package main

import (
	"dagger/file-sync/internal/dagger"
)

// buildContainer creates a container with rsync and SSH client configured for file sync
func (m *FileSync) buildContainer(source *dagger.Directory) *dagger.Container {
	container := dag.Container().
		From("alpine:3.20").
		WithMountedDirectory("/workspace", source).
		WithWorkdir("/workspace")

	// Install rsync and openssh-client
	container = container.WithExec([]string{
		"sh", "-c",
		"apk add --no-cache rsync openssh-client",
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
		Permissions: 0600,
	})

	return container
}
