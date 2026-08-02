package main

import (
	"context"
	"time"

	"dagger/terraform/internal/dagger"
)

// InitArtifact exécute `tofu init -backend=false` (télécharge providers + modules
// distants) et retourne le répertoire de travail initialisé, contenant
// .terraform/modules. Destiné à être scanné par security/trivy ScanConfig :
// contrairement à un scan du seul code wrapper, Trivy voit alors le contenu réel
// des modules distants via .terraform/modules/modules.json.
//
// Le sous-répertoire .terraform/providers (lourd) est retiré de l'artefact :
// inutile à `trivy config`, il alourdirait l'export.
func (m *Terraform) InitArtifact(
	ctx context.Context,
	// Répertoire contenant le code Terraform/OpenTofu
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
) (*dagger.Directory, error) {
	container := m.buildContainer(source, subpath)

	container, err := m.injectVariables(ctx, container)
	if err != nil {
		return nil, err
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"tofu", "init", "-backend=false"})

	return container.Directory(".").WithoutDirectory(".terraform/providers"), nil
}
