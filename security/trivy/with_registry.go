package main

import (
	"dagger/trivy/internal/dagger"
)

// WithRegistry configure l'authentification registry pour le scan d'images privées.
//
// Trivy s'authentifie via les variables TRIVY_USERNAME / TRIVY_PASSWORD lorsqu'il pull
// l'image à scanner. Utiliser env:VAR_NAME pour les identifiants.
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
