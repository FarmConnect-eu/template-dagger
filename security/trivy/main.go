// Package main provides a Dagger module for scanning IaC with Trivy.
//
// Trivy détecte nativement les misconfigurations dans un plan Terraform/OpenTofu
// exporté en JSON (`tofu show -json`). Ce module est conçu pour consommer
// l'artefact produit par le module iac/terraform (PlanArtifact).
package main

import (
	"context"

	"dagger/trivy/internal/dagger"
)

const defaultTrivyVersion = "0.70.0"

// defaultSeverity : niveaux de sévérité bloquants par défaut.
const defaultSeverity = "HIGH,CRITICAL"

type Trivy struct {
	// Version de l'image aquasec/trivy
	Version string
	// Sévérités à signaler (ex: "HIGH,CRITICAL")
	Severity string
	// Authentification registry (pour le scan d'images privées)
	RegistryHost     string
	RegistryUsername string
	RegistryPassword *dagger.Secret
}

func New() *Trivy {
	return &Trivy{
		Version:  defaultTrivyVersion,
		Severity: defaultSeverity,
	}
}

func (m *Trivy) Test(ctx context.Context) (string, error) {
	return m.buildContainer().
		WithExec([]string{"trivy", "--version"}).
		Stdout(ctx)
}
