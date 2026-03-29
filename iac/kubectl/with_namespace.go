package main

// WithNamespace sets the target Kubernetes namespace
func (m *Kubectl) WithNamespace(
	// Kubernetes namespace
	namespace string,
) *Kubectl {
	return &Kubectl{
		Kubeconfig:     m.Kubeconfig,
		Namespace:      namespace,
		KubectlVersion: m.KubectlVersion,
	}
}
