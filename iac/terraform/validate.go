package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Validate validates la configuration Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform validate`.
//
func (m *Terraform) Validate(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
) (string, error) {
	
	source, err := m.configureBackend(ctx, source, subpath)
	if err != nil {
		return "", err
	}

	
	container := m.buildContainer(source, subpath)

	
	container = container.WithExec([]string{"terraform", "init", "-backend=false"})

	
	container = container.WithExec([]string{"terraform", "validate"})

	
	return container.Stdout(ctx)
}
