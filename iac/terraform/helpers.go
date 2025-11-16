package main

import (
	"context"
	"dagger/terraform/internal/dagger"
	"fmt"
	"strings"
)

// resolveVariableValue resolves env:// or literal values
// Note: env:// prefix is stripped but not resolved - caller must pass resolved values
func (m *Terraform) resolveVariableValue(ctx context.Context, value string, isSecret bool) (any, error) {
	// Strip env:// prefix if present (caller must resolve environment variables)
	value = strings.TrimPrefix(value, "env://")

	if isSecret {
		return dag.SetSecret("literal", value), nil
	}
	return value, nil
}

// buildContainer creates Terraform container with source mounted
func (m *Terraform) buildContainer(
	source *dagger.Directory,
	workdir string,
) *dagger.Container {
	terraformVersion := m.TerraformVersion
	if terraformVersion == "" {
		terraformVersion = "1.9"
	}

	return dag.Container().
		From(fmt.Sprintf("hashicorp/terraform:%s", terraformVersion)).
		WithDirectory("/work", source).
		WithWorkdir(fmt.Sprintf("/work/%s", workdir))
}

// injectVariables injects variables as environment variables
func (m *Terraform) injectVariables(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	for _, v := range m.Variables {
		varName := v.Key
		if v.TfVar {
			varName = "TF_VAR_" + v.Key
		}

		resolvedValue, err := m.resolveVariableValue(ctx, v.Value, v.IsSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve variable %s: %w", v.Key, err)
		}

		if v.IsSecret {
			secret, ok := resolvedValue.(*dagger.Secret)
			if !ok {
				return nil, fmt.Errorf("expected secret for variable %s, got %T", v.Key, resolvedValue)
			}
			container = container.WithSecretVariable(varName, secret)
		} else {
			value, ok := resolvedValue.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for variable %s, got %T", v.Key, resolvedValue)
			}
			container = container.WithEnvVariable(varName, value)
		}
	}

	return container, nil
}

// configureBackend generates backend.tf if state is configured
func (m *Terraform) configureBackend(
	ctx context.Context,
	source *dagger.Directory,
	workdir string,
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
		return nil, fmt.Errorf("unsupported backend type: %s", m.State.Backend)
	}

	backendFile := dag.Directory().WithNewFile("backend.tf", backendConfig)
	source = source.WithDirectory(workdir, backendFile, dagger.DirectoryWithDirectoryOpts{
		Merge: true,
	})

	return source, nil
}
