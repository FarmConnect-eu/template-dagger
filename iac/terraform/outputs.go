package main

import (
	"context"
	"time"

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

	
	container, err = m.injectVariables(ctx, container, subpath)
	if err != nil {
		return "", err
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"tofu", "init"})

	args := []string{"tofu", "output"}
	if asJson {
		args = append(args, "-json")
	}
	if outputName != "" {
		args = append(args, outputName)
	}

	
	container = container.WithExec(args)


	return container.Stdout(ctx)
}

// SensitiveOutput retrieves a single Terraform output value, including sensitive ones
//
// Cette fonction utilise `terraform output -raw` pour récupérer la valeur brute
// d'un output, y compris les outputs marqués comme sensitive.
func (m *Terraform) SensitiveOutput(
	ctx context.Context,
	// Répertoire contenant le code Terraform
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Nom de l'output à récupérer (requis)
	outputName string,
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
		WithExec([]string{"tofu", "init"}).
		WithExec([]string{"tofu", "output", "-raw", outputName})

	return container.Stdout(ctx)
}
