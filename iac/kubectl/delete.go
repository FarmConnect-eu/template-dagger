package main

import (
	"context"
	"time"

	"dagger/kubectl/internal/dagger"
)

// Delete removes Kubernetes resources defined in the manifests
func (m *Kubectl) Delete(
	ctx context.Context,
	// Directory containing Kubernetes manifests
	source *dagger.Directory,
	// Subpath within source containing manifests
	// +optional
	// +default="."
	subpath string,
	// Ignore errors for resources that don't exist
	// +optional
	// +default=true
	ignoreNotFound bool,
) (string, error) {
	container := m.buildContainer(source, subpath)

	args := m.buildArgs([]string{"kubectl", "delete", "-f", ".", "--recursive"})

	if ignoreNotFound {
		args = append(args, "--ignore-not-found")
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args)

	return container.Stdout(ctx)
}
