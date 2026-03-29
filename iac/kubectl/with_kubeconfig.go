package main

import (
	"dagger/kubectl/internal/dagger"
)

// WithKubeconfig sets the kubeconfig for cluster authentication
func (m *Kubectl) WithKubeconfig(
	// Kubeconfig content as a secret
	kubeconfig *dagger.Secret,
) *Kubectl {
	return &Kubectl{
		Kubeconfig:     kubeconfig,
		Namespace:      m.Namespace,
		KubectlVersion: m.KubectlVersion,
	}
}
