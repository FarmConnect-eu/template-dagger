package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Validate valide la configuration Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform validate`.
//
// Exemple:
//
//	dagger call \
//	  validate --source ./terraform
func (m *Terraform) Validate(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
) (string, error) {
	// Configurer le backend si un état est configuré
	source, err := m.configureBackend(ctx, source)
	if err != nil {
		return "", err
	}

	// Construire le conteneur de base
	container := m.buildContainer(source)

	// Exécuter terraform init
	container = container.WithExec([]string{"terraform", "init", "-backend=false"})

	// Exécuter terraform validate
	container = container.WithExec([]string{"terraform", "validate"})

	// Retourner le résultat
	return container.Stdout(ctx)
}
