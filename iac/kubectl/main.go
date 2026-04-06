// Package main provides a Dagger module for deploying Kubernetes manifests via kubectl.
package main

import (
	"context"

	"dagger/kubectl/internal/dagger"
)

// Kubectl module for applying Kubernetes manifests in containers
type Kubectl struct {
	// Kubeconfig for cluster authentication
	Kubeconfig *dagger.Secret

	// Target namespace (default: "default")
	Namespace string
}

// New creates a new Kubectl instance with defaults
func New() *Kubectl {
	return &Kubectl{
		Kubeconfig: nil,
		Namespace:  "default",
	}
}

// Test verifies the module loads correctly
func (m *Kubectl) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "Kubectl module loaded successfully"}).
		Stdout(ctx)
}
