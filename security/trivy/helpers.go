package main

import (
	"context"
	"fmt"

	"dagger/trivy/internal/dagger"
)

// buildContainer prépare un conteneur Trivy à la version configurée.
//
// La base de vulnérabilités/checks est persistée dans un cache volume Dagger
// (`/root/.cache/trivy`) pour éviter un re-téléchargement à chaque scan.
func (m *Trivy) buildContainer() *dagger.Container {
	version := m.Version
	if version == "" {
		version = defaultTrivyVersion
	}
	return dag.Container().
		From(fmt.Sprintf("aquasec/trivy:%s", version)).
		WithMountedCache("/root/.cache/trivy", dag.CacheVolume("trivy-cache")).
		WithEnvVariable("TRIVY_NO_PROGRESS", "true")
}

// imageContainer prépare un conteneur Trivy avec, le cas échéant, les
// identifiants registry pour scanner une image privée (Trivy lit
// TRIVY_USERNAME / TRIVY_PASSWORD ; TRIVY_AUTH_URL restreint au registry visé).
func (m *Trivy) imageContainer() *dagger.Container {
	container := m.buildContainer()
	if m.RegistryUsername != "" && m.RegistryPassword != nil {
		container = container.
			WithEnvVariable("TRIVY_USERNAME", m.RegistryUsername).
			WithSecretVariable("TRIVY_PASSWORD", m.RegistryPassword)
		if m.RegistryHost != "" {
			container = container.WithEnvVariable("TRIVY_AUTH_URL", m.RegistryHost)
		}
	}
	return container
}

// reportExtension mappe un format Trivy vers une extension de fichier de rapport.
func reportExtension(format string) string {
	switch format {
	case "sarif":
		return "sarif"
	case "json", "cyclonedx":
		return "json"
	default:
		return "txt"
	}
}

// scanArtifact exécute Trivy une seule fois (résultat JSON en `/tmp/result.json`)
// sans faire échouer l'exec (Expect=Any), puis convertit ce résultat — sans
// re-scanner — vers le format demandé (`report.<ext>`) ET une table lisible
// (`report.txt`). Renvoie un répertoire contenant ces deux rapports plus un
// fichier `findings` valant "1" si des findings de la sévérité configurée
// existent, "0" sinon.
//
// La fonction ne retourne jamais d'erreur sur findings : les rapports sont ainsi
// toujours exportés (même avec findings), et le job CI décide de bloquer en
// lisant `findings`. Un seul scan Trivy produit à la fois l'artefact machine
// (ex: SARIF) et la table humaine pour le log.
//
// `args` doit lancer Trivy avec `--format json --output /tmp/result.json` et
// `--exit-code 1` (le code de sortie porte le verdict).
func scanArtifact(ctx context.Context, base *dagger.Container, args []string, format string) (*dagger.Directory, error) {
	executed := base.WithExec(args, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})

	code, err := executed.ExitCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("trivy scan: %w", err)
	}

	findings := "0"
	if code != 0 {
		findings = "1"
	}

	// Conversions à partir du JSON déjà produit (pas de re-scan).
	out := executed.
		WithExec([]string{"mkdir", "-p", "/out"}).
		WithExec([]string{"trivy", "convert", "--format", format, "--output", "/out/report." + reportExtension(format), "/tmp/result.json"})

	// Table humaine pour le log CI (inutile si le format demandé est déjà table).
	if format != "table" {
		out = out.WithExec([]string{"trivy", "convert", "--format", "table", "--output", "/out/report.txt", "/tmp/result.json"})
	}

	return out.
		WithNewFile("/out/findings", findings).
		Directory("/out"), nil
}
