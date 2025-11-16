package main

import (
	"context"
	"dagger/ansible/internal/dagger"
)

// ListInventory shows the Ansible inventory
//
// This function lists all hosts in the inventory.
// No variables needed for listing inventory.
//
// Example usage:
//
//	dagger call list-inventory --source=../infra-nfs
func (m *Ansible) ListInventory(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Inventory path (default: inventory/hosts.yml)
	// +optional
	// +default="inventory/hosts.yml"
	inventory string,
	// Working directory (default: ansible)
	// +optional
	// +default="ansible"
	workdir string,
) (string, error) {
	if inventory == "" {
		inventory = "inventory/hosts.yml"
	}
	if workdir == "" {
		workdir = "ansible"
	}

	// Create container with Ansible using the buildContainer helper
	container := m.buildContainer(source, workdir).
		WithExec([]string{"ansible-inventory", "-i", inventory, "--list"})

	return container.Stdout(ctx)
}
