package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger/terraform/internal/dagger"
)

func (m *Terraform) buildContainer(
	source *dagger.Directory,
	// +optional
	// +default="."
	subpath string,
) *dagger.Container {
	tofuVersion := m.TerraformVersion
	if tofuVersion == "" {
		tofuVersion = "1.10.6"
	}

	if subpath == "" {
		subpath = "."
	}

	return dag.Container().
		From(fmt.Sprintf("ghcr.io/opentofu/opentofu:%s", tofuVersion)).
		WithDirectory("/work", source).
		WithWorkdir(fmt.Sprintf("/work/%s", subpath))
}

// injectVariables adds environment variables and mounts tfvars files.
// Tfvars files are named *.auto.tfvars for automatic loading by Terraform.
func (m *Terraform) injectVariables(
	ctx context.Context,
	container *dagger.Container,
	// +optional
	// +default="."
	subpath string,
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

	// Mount tfvars files with .auto.tfvars extension for automatic loading
	for i, file := range m.TfVarsFiles {
		filename := fmt.Sprintf("dagger-%d.auto.tfvars", i)
		container = container.WithFile(filename, file)
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
		endpointBlock := ""
		if m.State.Endpoint != "" {
			endpointBlock = fmt.Sprintf(`
    endpoints {
      s3 = "%s"
    }`, m.State.Endpoint)
		}
		backendConfig = fmt.Sprintf(`terraform {
  backend "s3" {
    bucket                      = "%s"
    key                         = "%s"
    region                      = "%s"
    skip_requesting_account_id  = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    use_path_style              = true%s
  }
}
`, m.State.Bucket, m.State.Key, m.State.Region, endpointBlock)

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
