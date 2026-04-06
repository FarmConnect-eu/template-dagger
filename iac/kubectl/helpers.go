package main

import (
	"fmt"

	"dagger/kubectl/internal/dagger"
)

// buildContainer creates the base kubectl container with kubeconfig mounted
func (m *Kubectl) buildContainer(
	source *dagger.Directory,
	// +optional
	// +default="."
	subpath string,
) *dagger.Container {
	version := m.KubectlVersion
	if version == "" {
		version = "latest"
	}

	if subpath == "" {
		subpath = "."
	}

	container := dag.Container().
		From(fmt.Sprintf("bitnami/kubectl:%s", version)).
		WithDirectory("/work", source).
		WithWorkdir(fmt.Sprintf("/work/%s", subpath))

	// Mount kubeconfig as secret
	if m.Kubeconfig != nil {
		container = container.
			WithMountedSecret("/tmp/kubeconfig", m.Kubeconfig).
			WithEnvVariable("KUBECONFIG", "/tmp/kubeconfig")
	}

	return container
}

// buildArgs creates the common kubectl arguments (namespace, etc.)
func (m *Kubectl) buildArgs(baseArgs []string) []string {
	if m.Namespace != "" && m.Namespace != "default" {
		baseArgs = append(baseArgs, "-n", m.Namespace)
	}
	return baseArgs
}
