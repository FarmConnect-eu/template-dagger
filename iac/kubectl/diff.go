package main

import (
	"context"
	"time"

	"dagger/kubectl/internal/dagger"
)

// Diff shows differences between live cluster state and manifests
// Equivalent of "terraform plan" — shows what would change without applying
func (m *Kubectl) Diff(
	ctx context.Context,
	// Directory containing Kubernetes manifests
	source *dagger.Directory,
	// Subpath within source containing manifests
	// +optional
	// +default="."
	subpath string,
) (string, error) {
	container := m.buildContainer(source, subpath)

	args := m.buildArgs([]string{"kubectl", "diff", "-f", ".", "--recursive"})

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args)

	return container.Stdout(ctx)
}
