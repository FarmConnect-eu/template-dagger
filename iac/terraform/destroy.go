package main

import (
	"context"
	"time"

	"dagger/terraform/internal/dagger"
)

// Destroy destroys l'infrastructure gérée par Terraform
//
// Cette fonction exécute `terraform init` suivi de `terraform destroy`.
// Les variables doivent être configurées au préalable via WithVariable().
//
func (m *Terraform) Destroy(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Options supplémentaires pour terraform destroy
	// +optional
	destroyArgs []string,
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

	args := []string{"tofu", "destroy", "-auto-approve"}
	if len(destroyArgs) > 0 {
		args = append(args, destroyArgs...)
	}

	container = container.WithExec(args)

	return container.Stdout(ctx)
}
