package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Output retrieves les outputs Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform output`.
// Les outputs sont retournés au format JSON par défaut.
//
func (m *Terraform) Output(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Nom d'un output spécifique (laisser vide pour tous)
	// +optional
	outputName string,
	// Format JSON
	// +optional
	// +default=true
	asJson bool,
) (string, error) {
	
	source, err := m.configureBackend(ctx, source, subpath)
	if err != nil {
		return "", err
	}

	
	container := m.buildContainer(source, subpath)

	
	container, err = m.injectVariables(ctx, container)
	if err != nil {
		return "", err
	}

	
	container = container.WithExec([]string{"terraform", "init"})

	
	args := []string{"terraform", "output"}
	if asJson {
		args = append(args, "-json")
	}
	if outputName != "" {
		args = append(args, outputName)
	}

	
	container = container.WithExec(args)

	
	return container.Stdout(ctx)
}
