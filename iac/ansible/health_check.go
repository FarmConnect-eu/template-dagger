package main

import (
	"context"
	"dagger/ansible/internal/dagger"
)

// HealthCheck runs an Ansible health check playbook
//
// This function runs a health check playbook to verify service status.
// Variables must be configured beforehand using WithVariable().
//
// Example usage:
//
//	dagger call \
//	  with-variable --key ANSIBLE_HOST_KEY_CHECKING --value False \
//	  health-check \
//	  --source=../infra-nfs \
//	  --ssh-private-key=env:SSH_PRIVATE_KEY
func (m *Ansible) HealthCheck(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// SSH private key for connecting to hosts
	sshPrivateKey *dagger.Secret,
	// Inventory path (default: inventory/hosts.yml)
	// +optional
	// +default="inventory/hosts.yml"
	inventory string,
	// Limit to specific hosts
	// +optional
	limit string,
) (string, error) {
	if inventory == "" {
		inventory = "inventory/hosts.yml"
	}

	return m.AnsiblePlaybook(
		ctx,
		source,
		nil, // playbookSource (use source for both roles and playbooks)
		"ansible",
		"playbooks/health-check.yml",
		inventory,
		sshPrivateKey,
		"",
		"",
		limit,
		1, // verbose level 1
		false,
	)
}
