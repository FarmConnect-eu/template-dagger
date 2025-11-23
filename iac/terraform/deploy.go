package main

import (
	"context"

	"dagger/terraform/internal/dagger"
)

// Apply applies infrastructure changes Terraform à l'infrastructure
//
// Cette fonction exécute `terraform init` suivi de `terraform apply`.
// Les variables doivent être configurées au préalable via WithVariable().
//
func (m *Terraform) Apply(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Appliquer automatiquement sans confirmation
	// +optional
	// +default=false
	autoApprove bool,
	// Options supplémentaires pour terraform apply
	// +optional
	applyArgs []string,
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

	
	args := []string{"terraform", "apply"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	if len(applyArgs) > 0 {
		args = append(args, applyArgs...)
	}

	
	container = container.WithExec(args)

	
	return container.Stdout(ctx)
}
