package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Format formate les fichiers Terraform
//
// Cette fonction exécute `terraform fmt` pour formater le code.
//
// Exemple:
//
//	dagger call \
//	  format --source ./terraform
func (m *Terraform) Format(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Vérifier seulement si les fichiers sont formatés (sans modifier)
	// +optional
	// +default=false
	check bool,
	// Formatter récursivement les sous-répertoires
	// +optional
	// +default=true
	recursive bool,
) (string, error) {
	// Construire le conteneur de base
	container := m.buildContainer(source)

	// Construire la commande fmt
	args := []string{"terraform", "fmt"}
	if check {
		args = append(args, "-check")
	}
	if recursive {
		args = append(args, "-recursive")
	}

	// Exécuter terraform fmt
	container = container.WithExec(args)

	// Retourner le résultat
	return container.Stdout(ctx)
}
