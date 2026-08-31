package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
// identifiants registry pour scanner une image privée.
//
// On NE passe PAS par TRIVY_USERNAME / TRIVY_PASSWORD : Trivy interprète ces
// deux variables comme des LISTES séparées par des virgules et impose un
// appariement exact. Un identifiant contenant une virgule — cas courant des
// tokens de comptes robot Harbor — produit alors :
//
//	FATAL flag error: unable to convert flags to options:
//	the number of usernames and passwords must match
//
// et aucun échappement n'est prévu côté Trivy. On monte donc un
// `config.json` docker, que Trivy lit via DOCKER_CONFIG : le contenu est
// opaque pour lui, la virgule ne pose plus aucun problème.
//
// Le fichier est monté via WithMountedSecret et non écrit par un WithExec :
// un `printf > config.json` inscrirait le token dans une couche de conteneur,
// donc dans le cache Dagger. Ici il ne persiste nulle part.
func (m *Trivy) imageContainer(ctx context.Context) (*dagger.Container, error) {
	container := m.buildContainer()

	if m.RegistryUsername == "" || m.RegistryPassword == nil {
		return container, nil
	}
	if m.RegistryHost == "" {
		return nil, fmt.Errorf("registry host missing: use WithRegistry() with a host when credentials are set")
	}

	password, err := m.RegistryPassword.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("read registry password: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(m.RegistryUsername + ":" + password))
	config, err := json.Marshal(map[string]any{
		"auths": map[string]any{
			m.RegistryHost: map[string]string{"auth": auth},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("build docker config: %w", err)
	}

	// Nom de secret déterministe et dérivé des identifiants : deux registries
	// distincts dans une même session ne se marchent pas dessus.
	name := fmt.Sprintf("trivy-docker-config-%x",
		sha256.Sum256([]byte(m.RegistryHost+"\x00"+m.RegistryUsername)))

	return container.
		WithMountedSecret("/root/.docker/config.json", dag.SetSecret(name, string(config))).
		WithEnvVariable("DOCKER_CONFIG", "/root/.docker"), nil
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
