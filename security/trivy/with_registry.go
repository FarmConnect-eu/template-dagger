package main

import (
	"dagger/trivy/internal/dagger"
)

// WithRegistry configure l'authentification registry pour le scan d'images privées.
//
// Les identifiants sont transmis à Trivy via un `config.json` docker monté en secret
// (cf. imageContainer), et NON via TRIVY_USERNAME / TRIVY_PASSWORD : Trivy lit ces
// variables comme des listes séparées par des virgules, ce qui casse dès qu'un token
// en contient une. Utiliser env:VAR_NAME pour le mot de passe.
func (m *Trivy) WithRegistry(
	// Hôte du registry (ex: "dordogne.azurecr.io")
	host string,
	// Nom d'utilisateur du registry
	username string,
	// Mot de passe / token du registry
	password *dagger.Secret,
) *Trivy {
	return &Trivy{
		Version:          m.Version,
		Severity:         m.Severity,
		RegistryHost:     host,
		RegistryUsername: username,
		RegistryPassword: password,
	}
}
