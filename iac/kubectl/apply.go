package main

import (
	"context"
	"time"

	"dagger/kubectl/internal/dagger"
)

// Apply applies Kubernetes manifests from the source directory
func (m *Kubectl) Apply(
	ctx context.Context,
	// Directory containing Kubernetes manifests
	source *dagger.Directory,
	// Subpath within source containing manifests
	// +optional
	// +default="."
	subpath string,
	// Prune resources not in the manifest set
	// +optional
	// +default=false
	prune bool,
	// Run in server-side dry-run mode (no changes)
	// +optional
	// +default=false
	dryRun bool,
) (string, error) {
	container := m.buildContainer(source, subpath)

	args := m.buildArgs([]string{"kubectl", "apply", "-f", ".", "--recursive"})

	if prune {
		args = append(args, "--prune")
	}
	if dryRun {
		args = append(args, "--dry-run=server")
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args)

	return container.Stdout(ctx)
}
