package main

import (
	"context"
	"dagger/ansible/internal/dagger"
)

// ConfigureWithAnsible runs an Ansible playbook to configure infrastructure
//
// This function executes Ansible playbooks after infrastructure deployment.
// It requires SSH access to the target hosts.
// Variables must be configured beforehand using WithVariable().
//
// Example usage:
//
//	dagger call \
//	  with-variable --key ANSIBLE_HOST_KEY_CHECKING --value False \
//	  configure-with-ansible \
//	  --source=../infra-nfs \
//	  --playbook=playbooks/deploy-nfs.yml \
//	  --inventory=inventory/hosts.yml \
//	  --ssh-private-key=env:SSH_PRIVATE_KEY
func (m *Ansible) ConfigureWithAnsible(
	ctx context.Context,
	// Infrastructure repository directory
	source *dagger.Directory,
	// Path to playbook relative to source (e.g., "playbooks/deploy-nfs.yml")
	playbook string,
	// Path to inventory relative to source (default: "inventory/hosts.yml")
	// +optional
	// +default="inventory/hosts.yml"
	inventory string,
	// SSH private key for connecting to hosts
	sshPrivateKey *dagger.Secret,
	// Working directory relative to source (default: ansible)
	// +optional
	// +default="ansible"
	workdir string,
	// Extra variables to pass to ansible-playbook (JSON format)
	// +optional
	extraVars string,
	// Ansible tags to run specific tasks
	// +optional
	tags string,
	// Limit execution to specific hosts
	// +optional
	limit string,
	// Verbose output level (0-4, where 4 is max verbosity)
	// +optional
	// +default=0
	verbose int,
	// Check mode (dry-run without making changes)
	// +optional
	// +default=false
	checkMode bool,
) (string, error) {
	if workdir == "" {
		workdir = "ansible"
	}
	if inventory == "" {
		inventory = "inventory/hosts.yml"
	}

	return m.AnsiblePlaybook(
		ctx,
		source,
		nil, // playbookSource (use source for both roles and playbooks)
		workdir,
		playbook,
		inventory,
		sshPrivateKey,
		extraVars,
		tags,
		limit,
		verbose,
		checkMode,
	)
}
