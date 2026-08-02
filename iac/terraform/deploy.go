package main

import (
	"context"
	"time"

	"dagger/terraform/internal/dagger"
)

// Apply applique les changements d'infrastructure avec OpenTofu
//
// Sans planFile : exécute `tofu init` puis `tofu apply -auto-approve` (recalcule le plan).
// Avec planFile : exécute `tofu init` puis `tofu apply -auto-approve tfplan.bin`, appliquant
// exactement le plan produit par PlanArtifact (pas de recalcul). Si l'état a changé depuis
// la génération du plan, OpenTofu refuse le plan obsolète (« saved plan is stale »).
//
// Les credentials providers restent nécessaires même avec un plan sauvegardé : ils ne sont
// pas stockés dans le plan. Les variables doivent donc toujours être configurées au préalable.
func (m *Terraform) Apply(
	ctx context.Context,
	// Répertoire contenant le code Terraform/OpenTofu
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Plan sauvegardé (tfplan.bin) à appliquer tel quel, généré par PlanArtifact
	// +optional
	planFile *dagger.File,
	// Options supplémentaires pour tofu apply
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

	if planFile != nil {
		container = container.WithFile(planFileName, planFile)
	}

	container = container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"tofu", "init"})

	args := []string{"tofu", "apply", "-auto-approve"}
	if len(applyArgs) > 0 {
		args = append(args, applyArgs...)
	}
	if planFile != nil {
		// Le fichier de plan doit être le dernier argument positionnel.
		args = append(args, planFileName)
	}

	container = container.WithExec(args)

	return container.Stdout(ctx)
}
