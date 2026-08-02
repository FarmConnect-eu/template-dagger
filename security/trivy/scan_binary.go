package main

import (
	"context"

	"dagger/trivy/internal/dagger"
)

// ScanBinary scanne un binaire compilé (`trivy rootfs`).
//
// Pour un binaire Go, Trivy lit les métadonnées de build embarquées par
// `go build` et détecte les CVE sur les dépendances. Le binaire produit
// par languages/go (Build) peut être chaîné directement en entrée.
//
// Limitation : ne fonctionne pas sur les binaires compressés avec UPX
// (les métadonnées embarquées ne sont plus lisibles).
func (m *Trivy) ScanBinary(
	ctx context.Context,
	// Binaire à scanner (typiquement issu de languages/go Build)
	binary *dagger.File,
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
		"trivy", "rootfs",
		"--severity", m.Severity,
		"--format", format,
	}
	if failOnFindings {
		args = append(args, "--exit-code", "1")
	}
	args = append(args, "/scan/binary")

	return m.buildContainer().
		WithMountedFile("/scan/binary", binary).
		WithExec(args).
		Stdout(ctx)
}
