package main

import (
	"context"
	"fmt"

	"dagger/terraform/internal/dagger"
)

// buildContainer crée un conteneur Terraform avec le code source monté
func (m *Terraform) buildContainer(
	source *dagger.Directory,
) *dagger.Container {
	terraformVersion := m.TerraformVersion
	if terraformVersion == "" {
		terraformVersion = "1.9.8"
	}

	return dag.Container().
		From(fmt.Sprintf("hashicorp/terraform:%s", terraformVersion)).
		WithDirectory("/work", source).
		WithWorkdir("/work")
}

// injectVariables injecte les variables accumulées dans le conteneur
// Les variables marquées comme secrets sont injectées via WithSecretVariable
// Les autres via WithEnvVariable
// Le préfixe env:// est géré nativement par Dagger (pas besoin de le stripper)
func (m *Terraform) injectVariables(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	for _, v := range m.Variables {
		varName := v.Key
		if v.TfVar {
			varName = "TF_VAR_" + v.Key
		}

		if v.IsSecret {
			// Pour les secrets, on utilise SetSecret avec la valeur
			// Dagger gère automatiquement env://, file://, etc.
			secret := dag.SetSecret(v.Key, v.Value)
			container = container.WithSecretVariable(varName, secret)
		} else {
			// Pour les variables non-secrètes, on injecte directement
			// Note: Si la valeur contient env://, Dagger le résoudra automatiquement
			container = container.WithEnvVariable(varName, v.Value)
		}
	}

	return container, nil
}

// configureBackend génère dynamiquement le fichier backend.tf si un état est configuré
// Supporte : s3, gcs, azurerm, local
func (m *Terraform) configureBackend(
	ctx context.Context,
	source *dagger.Directory,
) (*dagger.Directory, error) {
	if m.State == nil {
		return source, nil
	}

	var backendConfig string

	switch m.State.Backend {
	case "s3":
		backendConfig = fmt.Sprintf(`terraform {
  backend "s3" {
    bucket = "%s"
    key    = "%s"
    region = "%s"
  }
}
`, m.State.Bucket, m.State.Key, m.State.Region)

	case "gcs":
		backendConfig = fmt.Sprintf(`terraform {
  backend "gcs" {
    bucket = "%s"
    prefix = "%s"
  }
}
`, m.State.Bucket, m.State.Key)

	case "azurerm":
		backendConfig = fmt.Sprintf(`terraform {
  backend "azurerm" {
    container_name = "%s"
    key           = "%s"
  }
}
`, m.State.Bucket, m.State.Key)

	case "local":
		backendConfig = fmt.Sprintf(`terraform {
  backend "local" {
    path = "%s"
  }
}
`, m.State.Key)

	default:
		return nil, fmt.Errorf("unsupported backend type: %s (supported: s3, gcs, azurerm, local)", m.State.Backend)
	}

	// Créer le fichier backend.tf à la racine du source
	source = source.WithNewFile("backend.tf", backendConfig)

	return source, nil
}
