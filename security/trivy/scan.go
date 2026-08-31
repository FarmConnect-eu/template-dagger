package main

import (
	"context"
	"time"

	"dagger/trivy/internal/dagger"
)

// ScanImage scanne les vulnérabilités d'une image Docker (`trivy image`).
//
// Pour une image d'un registry privé, configurer les identifiants au préalable via WithRegistry() :
// Trivy pull l'image lui-même, en s'authentifiant via un config.json docker monté (cf. imageContainer).
func (m *Trivy) ScanImage(
	ctx context.Context,
	// Référence complète de l'image (ex: "dordogne.azurecr.io/ad-frontend:v5")
	image string,
	// Ignorer les vulnérabilités sans correctif disponible
	// +optional
	// +default=true
	ignoreUnfixed bool,
	// Faire échouer le scan si des vulnérabilités sont détectées
	// +optional
	// +default=true
	failOnFindings bool,
	// Format de sortie: table, json, sarif
	// +optional
	// +default="table"
	format string,
) (string, error) {
	container, err := m.imageContainer(ctx)
	if err != nil {
		return "", err
	}

	args := []string{
		"trivy", "image",
		"--severity", m.Severity,
		"--format", format,
	}
	if ignoreUnfixed {
		args = append(args, "--ignore-unfixed")
	}
	if failOnFindings {
		args = append(args, "--exit-code", "1")
	}
	args = append(args, image)

	// Le tag d'image est mutable (peut être re-poussé) → empêcher la mise en cache du résultat.
	return container.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec(args).
		Stdout(ctx)
}

// ScanImageArtifact scanne une image (`trivy image`) en un seul passage et renvoie
// un répertoire contenant le rapport (`report.<ext>`) et un fichier `findings`
// ("1" si vulnérabilités de la sévérité configurée, "0" sinon).
//
// Comme ScanFsArtifact, l'appel n'échoue jamais sur findings : le rapport est
// toujours produit (à exporter via `export --path`), et la barrière CI se décide
// en lisant `findings`. Fusionne le couple SARIF + gate en un seul scan.
// Le SBOM (inventaire complet) reste une passe distincte (portée différente).
func (m *Trivy) ScanImageArtifact(
	ctx context.Context,
	// Référence complète de l'image (idéalement par digest : registre/nom@sha256:...)
	image string,
	// Ignorer les vulnérabilités sans correctif disponible
	// +optional
	// +default=true
	ignoreUnfixed bool,
	// Format du rapport: table, json, sarif (défaut: sarif)
	// +optional
	// +default="sarif"
	format string,
) (*dagger.Directory, error) {
	if format == "" {
		format = "sarif"
	}

	args := []string{
		"trivy", "image",
		"--severity", m.Severity,
	}
	if ignoreUnfixed {
		args = append(args, "--ignore-unfixed")
	}
	// Résultat JSON réutilisé par scanArtifact (convert -> format demandé + table).
	args = append(args, "--format", "json", "--output", "/tmp/result.json", "--exit-code", "1", image)

	// Image potentiellement mutable (tag) → empêcher la mise en cache du résultat.
	container, err := m.imageContainer(ctx)
	if err != nil {
		return nil, err
	}
	base := container.WithEnvVariable("CACHEBUSTER", time.Now().String())

	return scanArtifact(ctx, base, args, format)
}

// ScanPlan scanne un plan Terraform/OpenTofu exporté en JSON (`tofu show -json`).
//
// Trivy détecte automatiquement le format « Terraform Plan JSON » via `trivy config`.
// Avec failOnFindings=true, la commande retourne un code non nul si une misconfiguration
// d'une sévérité configurée est trouvée — ce qui fait échouer le job CI.
func (m *Trivy) ScanPlan(
	ctx context.Context,
	// Plan au format JSON (tfplan.json) produit par iac/terraform PlanArtifact
	plan *dagger.File,
	// Faire échouer le scan si des findings sont détectés
	// +optional
	// +default=true
	failOnFindings bool,
	// Format de sortie: table, json, sarif
	// +optional
	// +default="table"
	format string,
) (string, error) {
	args := []string{
		"trivy", "config",
		"--severity", m.Severity,
		"--format", format,
	}
	if failOnFindings {
		args = append(args, "--exit-code", "1")
	}
	args = append(args, "tfplan.json")

	return m.buildContainer().
		WithMountedFile("/scan/tfplan.json", plan).
		WithWorkdir("/scan").
		WithExec(args).
		Stdout(ctx)
}

// ScanConfig scanne du code IaC brut (HCL, Dockerfile, Kubernetes...) dans un répertoire.
func (m *Trivy) ScanConfig(
	ctx context.Context,
	// Répertoire contenant le code à scanner
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Faire échouer le scan si des findings sont détectés
	// +optional
	// +default=true
	failOnFindings bool,
	// Format de sortie: table, json, sarif
	// +optional
	// +default="table"
	format string,
) (string, error) {
	if subpath == "" {
		subpath = "."
	}

	args := []string{
		"trivy", "config",
		"--severity", m.Severity,
		"--format", format,
	}
	if failOnFindings {
		args = append(args, "--exit-code", "1")
	}
	args = append(args, subpath)

	return m.buildContainer().
		WithMountedDirectory("/scan", source).
		WithWorkdir("/scan").
		WithExec(args).
		Stdout(ctx)
}
