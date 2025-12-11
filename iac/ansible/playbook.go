package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/ansible/internal/dagger"
)

// RunPlaybook executes an Ansible playbook with the configured inventory, variables, and secrets.
func (m *Ansible) RunPlaybook(
	ctx context.Context,
	// Playbook file (e.g., site.yml)
	playbook *dagger.File,
	// Source directory (contains roles/, group_vars/, etc.)
	source *dagger.Directory,
	// SSH private key for host connections (supports file:, env:)
	sshPrivateKey *dagger.Secret,
	// Check mode (dry-run, no changes)
	// +optional
	// +default=false
	checkMode bool,
	// Limit pattern to restrict execution to specific hosts
	// +optional
	limit string,
) (string, error) {
	if playbook == nil {
		return "", fmt.Errorf("playbook file is required")
	}

	if m.Inventory == nil {
		return "", fmt.Errorf("inventory is required: use WithInventory to set inventory file")
	}

	if sshPrivateKey == nil {
		return "", fmt.Errorf("ssh-private-key is required")
	}

	container := m.buildContainer(source)

	container, err := m.injectVariables(ctx, container)
	if err != nil {
		return "", fmt.Errorf("failed to inject variables: %w", err)
	}

	container = container.
		WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
		WithMountedSecret("/root/.ssh/id_rsa", sshPrivateKey, dagger.ContainerWithMountedSecretOpts{
			Mode:  0600,
			Owner: "root:root",
		})

	container = container.
		WithExec([]string{"sh", "-c", "ssh-keyscan -H github.com >> /root/.ssh/known_hosts || true"})

	container = container.WithMountedFile("/work/inventory", m.Inventory)
	container = container.WithMountedFile("/work/playbook.yml", playbook)

	args := []string{"/opt/ansible-venv/bin/ansible-playbook"}
	args = append(args, "-i", "/work/inventory")
	args = append(args, "/work/playbook.yml")

	args = append(args, "-v")

	if checkMode {
		args = append(args, "--check")
	}

	if len(m.Tags) > 0 {
		args = append(args, "--tags", strings.Join(m.Tags, ","))
	}

	if len(m.SkipTags) > 0 {
		args = append(args, "--skip-tags", strings.Join(m.SkipTags, ","))
	}

	if limit != "" {
		args = append(args, "--limit", limit)
	}

	if len(m.ExtraVars) > 0 {
		extraVarsJSON := m.buildExtraVars()
		args = append(args, "--extra-vars", extraVarsJSON)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute playbook: %w", err)
	}

	return output, nil
}
