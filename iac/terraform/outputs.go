package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Output récupère les outputs Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform output`.
// Les outputs sont retournés au format JSON par défaut.
//
// Exemple (tous les outputs):
//
//	dagger call \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  output --source ./terraform
//
// Exemple (output spécifique):
//
//	dagger call \
//	  with-state --backend s3 --bucket my-state --key terraform.tfstate --region us-east-1 \
//	  output --source ./terraform --output-name instance_ip
func (m *Terraform) Output(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Nom d'un output spécifique (laisser vide pour tous)
	// +optional
	outputName string,
	// Format JSON
	// +optional
	// +default=true
	asJson bool,
) (string, error) {
	// Configurer le backend si un état est configuré
	source, err := m.configureBackend(ctx, source)
	if err != nil {
		return "", err
	}

	// Construire le conteneur de base
	container := m.buildContainer(source)

	// Injecter les variables (pour l'accès au backend)
	container, err = m.injectVariables(ctx, container)
	if err != nil {
		return "", err
	}

	// Exécuter terraform init
	container = container.WithExec([]string{"terraform", "init"})

	// Construire la commande output
	args := []string{"terraform", "output"}
	if asJson {
		args = append(args, "-json")
	}
	if outputName != "" {
		args = append(args, outputName)
	}

	// Exécuter terraform output
	container = container.WithExec(args)

	// Retourner le résultat
	return container.Stdout(ctx)
}
