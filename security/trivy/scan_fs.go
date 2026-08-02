package main

import (
	"context"

	"dagger/trivy/internal/dagger"
)

// ScanFs scanne le filesystem d'un projet (`trivy fs`).
//
// Trivy auto-détecte les lockfiles présents : go.mod/go.sum, package-lock.json,
// requirements.txt, Cargo.lock, etc. Les scanners actifs sont contrôlés via
// le paramètre `scanners` (vide => défauts Trivy: vuln,secret).
func (m *Trivy) ScanFs(
	ctx context.Context,
	// Répertoire à scanner
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Liste de scanners séparés par virgule (ex: "vuln", "vuln,secret", "vuln,secret,misconfig").
	// Vide => défauts de Trivy (vuln,secret).
	// +optional
	scanners string,
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
		"trivy", "fs",
		"--severity", m.Severity,
		"--format", format,
	}
	if scanners != "" {
		args = append(args, "--scanners", scanners)
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

// ScanFsArtifact scanne le filesystem (`trivy fs`) en un seul passage et renvoie
// un répertoire contenant le rapport (`report.<ext>`) et un fichier `findings`
// ("1" si findings HIGH/CRITICAL, "0" sinon).
//
// Contrairement à ScanFs, l'appel n'échoue jamais sur findings : le rapport est
// toujours produit (à exporter via `export --path`), et la barrière CI se décide
// en lisant `findings`. Cela remplace le double scan (rapport + gate) par un seul.
func (m *Trivy) ScanFsArtifact(
	ctx context.Context,
	// Répertoire à scanner
	source *dagger.Directory,
	// Sous-chemin relatif dans source (défaut: ".")
	// +optional
	// +default="."
	subpath string,
	// Liste de scanners séparés par virgule (vide => défauts Trivy: vuln,secret)
	// +optional
	scanners string,
	// Format du rapport: table, json, sarif (défaut: sarif)
	// +optional
	// +default="sarif"
	format string,
) (*dagger.Directory, error) {
	if subpath == "" {
		subpath = "."
	}
	if format == "" {
		format = "sarif"
	}

	args := []string{
		"trivy", "fs",
		"--severity", m.Severity,
	}
	if scanners != "" {
		args = append(args, "--scanners", scanners)
	}
	// Résultat JSON réutilisé par scanArtifact (convert -> format demandé + table).
	args = append(args, "--format", "json", "--output", "/tmp/result.json", "--exit-code", "1", subpath)

	base := m.buildContainer().
		WithMountedDirectory("/scan", source).
		WithWorkdir("/scan")

	return scanArtifact(ctx, base, args, format)
}
