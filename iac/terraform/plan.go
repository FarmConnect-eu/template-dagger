package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Plan génère et affiche un plan d'exécution Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform plan`.
// Les variables doivent être configurées au préalable via WithVariable().
//
func (m *Terraform) Plan(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Utiliser -detailed-exitcode (0=no changes, 1=error, 2=changes)
	// +optional
	// +default=false
	detailedExitcode bool,
	// Options supplémentaires pour terraform plan
	// +optional
	planArgs []string,
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

	
	args := []string{"terraform", "plan"}
	if detailedExitcode {
		args = append(args, "-detailed-exitcode")
	}
	if len(planArgs) > 0 {
		args = append(args, planArgs...)
	}

	
	container = container.WithExec(args)

	
	return container.Stdout(ctx)
}
