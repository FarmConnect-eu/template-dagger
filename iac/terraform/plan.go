package main

import (
	"context"
	"fmt"

	"dagger/terraform/internal/dagger"
)

const planFileName = "tfplan.bin"

// Plan génère et affiche un plan d'exécution Terraform
//
// Cette fonction exécute `tofu init` suivi de `tofu plan`.
// Les variables doivent être configurées au préalable via WithVariable().
//
// Pour produire un plan réutilisable (consommé par Apply --plan-file) et scannable
// par Trivy, utiliser PlanArtifact qui retourne un répertoire d'artefacts.
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
	// Générer un plan de destruction (tofu plan -destroy)
	// +optional
	// +default=false
	destroy bool,
	// Options supplémentaires pour terraform plan
	// +optional
	planArgs []string,
) (string, error) {

	container, err := m.preparePlanContainer(ctx, source, subpath)
	if err != nil {
		return "", err
	}

	args := []string{"tofu", "plan"}
	if destroy {
		args = append(args, "-destroy")
	}
	if detailedExitcode {
		args = append(args, "-detailed-exitcode")
	}
	if len(planArgs) > 0 {
		args = append(args, planArgs...)
	}

	return container.WithExec(args).Stdout(ctx)
}

// PlanArtifact génère un plan d'exécution sauvegardé et réutilisable.
//
// Contrairement à Plan (qui n'affiche que la sortie texte), cette fonction
// exécute `tofu plan -out=tfplan.bin` puis retourne un répertoire contenant :
//   - tfplan.bin   : le plan binaire, consommable tel quel par Apply --plan-file
//   - tfplan.json  : `tofu show -json` du plan (pour le scan Trivy)
//   - plan.txt     : `tofu show` lisible pour les reviewers (valeurs sensibles masquées)
//   - changes      : "true" si des changements sont en attente, sinon "false"
//
// Les variables doivent être configurées au préalable via WithVariable()/WithSecret().
func (m *Terraform) PlanArtifact(
	ctx context.Context,
	// Répertoire contenant le code Terraform/OpenTofu
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Générer un plan de destruction (tofu plan -destroy)
	// +optional
	// +default=false
	destroy bool,
	// Options supplémentaires pour tofu plan
	// +optional
	planArgs []string,
) (*dagger.Directory, error) {

	container, err := m.preparePlanContainer(ctx, source, subpath)
	if err != nil {
		return nil, err
	}

	// tofu plan -detailed-exitcode -out=tfplan.bin
	// -detailed-exitcode: 0 = aucun changement, 1 = erreur, 2 = changements.
	args := []string{"tofu", "plan", "-detailed-exitcode", "-out=" + planFileName}
	if destroy {
		args = append(args, "-destroy")
	}
	if len(planArgs) > 0 {
		args = append(args, planArgs...)
	}

	// Expect=ANY pour récupérer le code de sortie 2 (changements) sans échec.
	planned := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: dagger.ReturnTypeAny,
	})

	exitCode, err := planned.ExitCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("plan execution failed: %w", err)
	}
	if exitCode == 1 {
		stderr, _ := planned.Stderr(ctx)
		return nil, fmt.Errorf("tofu plan failed: %s", stderr)
	}
	hasChanges := exitCode == 2

	// Sortie texte lisible (les valeurs sensibles sont masquées par tofu show).
	planText, err := planned.WithExec([]string{"tofu", "show", "-no-color", planFileName}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("tofu show failed: %w", err)
	}

	// Sortie JSON pour le scan Trivy.
	planJSON, err := planned.WithExec([]string{"tofu", "show", "-json", planFileName}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("tofu show -json failed: %w", err)
	}

	changes := "false"
	if hasChanges {
		changes = "true"
	}

	out := dag.Directory().
		WithFile(planFileName, planned.File(planFileName)).
		WithNewFile("tfplan.json", planJSON).
		WithNewFile("plan.txt", planText).
		WithNewFile("changes", changes)

	return out, nil
}
