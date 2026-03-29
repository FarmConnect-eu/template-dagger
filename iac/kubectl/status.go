package main

import (
	"context"
	"time"

	"dagger/kubectl/internal/dagger"
)

// Status shows all resources in the target namespace
func (m *Kubectl) Status(
	ctx context.Context,
	// Directory (needed for Dagger source mounting, can be empty)
	source *dagger.Directory,
) (string, error) {
	container := m.buildContainer(source, ".")

	args := m.buildArgs([]string{"kubectl", "get", "all", "-o", "wide"})

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args)

	return container.Stdout(ctx)
}

// Rollout checks or restarts a deployment rollout
func (m *Kubectl) Rollout(
	ctx context.Context,
	// Directory (needed for Dagger source mounting)
	source *dagger.Directory,
	// Deployment name to check/restart
	deployment string,
	// Restart the deployment instead of checking status
	// +optional
	// +default=false
	restart bool,
) (string, error) {
	container := m.buildContainer(source, ".")

	var args []string
	if restart {
		args = m.buildArgs([]string{"kubectl", "rollout", "restart", "deployment/" + deployment})
	} else {
		args = m.buildArgs([]string{"kubectl", "rollout", "status", "deployment/" + deployment})
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args)

	return container.Stdout(ctx)
}
