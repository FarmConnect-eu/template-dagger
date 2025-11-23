package main

import (
	"context"
	"fmt"

	"dagger/terraform/internal/dagger"
)

func (m *Terraform) buildContainer(
	source *dagger.Directory,
	// +optional
	// +default="."
	subpath string,
) *dagger.Container {
	terraformVersion := m.TerraformVersion
	if terraformVersion == "" {
		terraformVersion = "1.9.8"
	}

	if subpath == "" {
		subpath = "."
	}

	return dag.Container().
		From(fmt.Sprintf("hashicorp/terraform:%s", terraformVersion)).
		WithDirectory("/work", source).
		WithWorkdir(fmt.Sprintf("/work/%s", subpath))
}

func (m *Terraform) injectVariables(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	for _, v := range m.Variables {
		varName := v.Key
		if v.TfVar {
			varName = "TF_VAR_" + v.Key
		}

		if v.SecretValue != nil {
			container = container.WithSecretVariable(varName, v.SecretValue)
		} else {
			container = container.WithEnvVariable(varName, v.Value)
		}
	}

	return container, nil
}

func (m *Terraform) configureBackend(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="."
	subpath string,
) (*dagger.Directory, error) {
	if m.State == nil {
		return source, nil
	}

	if subpath == "" {
		subpath = "."
	}

	var backendConfig string

	switch m.State.Backend {
	case "s3":
		backendConfig = fmt.Sprintf(`terraform {
  backend "s3" {
    bucket                      = "%s"
    key                         = "%s"
    region                      = "%s"
    skip_requesting_account_id  = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    use_path_style              = true
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

	backendPath := fmt.Sprintf("%s/backend.tf", subpath)
	if subpath == "." {
		backendPath = "backend.tf"
	}
	source = source.WithNewFile(backendPath, backendConfig)

	return source, nil
}
