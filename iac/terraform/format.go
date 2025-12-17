package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Format formats les fichiers Terraform
//
// Cette fonction exécute `terraform fmt` pour formater le code.
//
func (m *Terraform) Format(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Vérifier seulement si les fichiers sont formatés (sans modifier)
	// +optional
	// +default=false
	check bool,
	// Formatter récursivement les sous-répertoires
	// +optional
	// +default=true
	recursive bool,
) (string, error) {
	
	container := m.buildContainer(source, subpath)

	
	args := []string{"tofu", "fmt"}
	if check {
		args = append(args, "-check")
	}
	if recursive {
		args = append(args, "-recursive")
	}

	
	container = container.WithExec(args)

	
	return container.Stdout(ctx)
}
