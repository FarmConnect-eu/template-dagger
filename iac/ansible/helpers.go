package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/ansible/internal/dagger"
)

// resolveVariableValue resolves env:// or literal values
// Note: env:// prefix is stripped but not resolved - caller must pass resolved values
func (m *Ansible) resolveVariableValue(ctx context.Context, value string, isSecret bool) (any, error) {
	// Strip env:// prefix if present (caller must resolve environment variables)
	value = strings.TrimPrefix(value, "env://")

	if isSecret {
		return dag.SetSecret("literal", value), nil
	}
	return value, nil
}

// injectVariables injects variables as environment variables
func (m *Ansible) injectVariables(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	for _, v := range m.Variables {
		resolvedValue, err := m.resolveVariableValue(ctx, v.Value, v.IsSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve variable %s: %w", v.Key, err)
		}

		if v.IsSecret {
			secret, ok := resolvedValue.(*dagger.Secret)
			if !ok {
				return nil, fmt.Errorf("expected secret for variable %s, got %T", v.Key, resolvedValue)
			}
			container = container.WithSecretVariable(v.Key, secret)
		} else {
			value, ok := resolvedValue.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for variable %s, got %T", v.Key, resolvedValue)
			}
			container = container.WithEnvVariable(v.Key, value)
		}
	}

	return container, nil
}

// buildContainer creates Ansible container on Ubuntu 22.04 with source mounted
func (m *Ansible) buildContainer(
	source *dagger.Directory,
	workdir string,
) *dagger.Container {
	ansibleVersion := m.AnsibleVersion
	if ansibleVersion == "" {
		ansibleVersion = "2.15"
	}

	return dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{
			"apt-get", "install", "-y",
			"software-properties-common",
			"python3",
			"python3-pip",
			"openssh-client",
			"sshpass",
			"git",
			"curl",
		}).
		WithExec([]string{"add-apt-repository", "--yes", "--update", "ppa:ansible/ansible"}).
		WithExec([]string{"apt-get", "install", "-y", "ansible=" + ansibleVersion + "*"}).
		WithExec([]string{"apt-get", "clean"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).
		WithDirectory("/work", source).
		WithWorkdir("/work/" + workdir)
}
