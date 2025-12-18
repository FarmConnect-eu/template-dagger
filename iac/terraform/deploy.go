package main

import (
	"context"
	"time"

	"dagger/terraform/internal/dagger"
)

// Apply applique les changements d'infrastructure avec OpenTofu
//
// Cette fonction exécute `tofu init` suivi de `tofu apply -auto-approve`.
// Les variables doivent être configurées au préalable via WithVariable().
//
func (m *Terraform) Apply(
	ctx context.Context,
	// Répertoire contenant le code Terraform/OpenTofu
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Options supplémentaires pour tofu apply
	// +optional
	applyArgs []string,
) (string, error) {

	source, err := m.configureBackend(ctx, source, subpath)
	if err != nil {
		return "", err
	}

	container := m.buildContainer(source, subpath)

	container, err = m.injectVariables(ctx, container, subpath)
	if err != nil {
		return "", err
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"tofu", "init"})

	args := []string{"tofu", "apply", "-auto-approve"}
	if len(applyArgs) > 0 {
		args = append(args, applyArgs...)
	}

	container = container.WithExec(args)

	return container.Stdout(ctx)
}
