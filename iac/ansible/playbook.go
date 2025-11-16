package main

import (
	"context"
	"dagger/ansible/internal/dagger"
)

// AnsiblePlaybook executes an Ansible playbook
// source = template-ansible (roles), playbookSource = infra-* (playbooks/inventory)
func (m *Ansible) AnsiblePlaybook(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	playbookSource *dagger.Directory,
	workdir string,
	playbook string,
	inventory string,
	sshPrivateKey *dagger.Secret,
	extraVars string,
	tags string,
	limit string,
	verbose int,
	checkMode bool,
) (string, error) {
	// Build ansible-playbook command
	args := []string{
		"ansible-playbook",
		"-i", inventory,
		playbook,
	}

	// Add verbose flags
	if verbose > 0 {
		verboseFlag := "-"
		for i := 0; i < verbose && i < 4; i++ {
			verboseFlag += "v"
		}
		args = append(args, verboseFlag)
	}

	// Add check mode
	if checkMode {
		args = append(args, "--check")
	}

	// Add tags
	if tags != "" {
		args = append(args, "--tags", tags)
	}

	// Add limit
	if limit != "" {
		args = append(args, "--limit", limit)
	}

	// Add extra vars
	if extraVars != "" {
		args = append(args, "--extra-vars", extraVars)
	}

	// Determine sources: if playbookSource is nil, use source for everything (backward compat)
	rolesSource := source
	workSource := source
	if playbookSource != nil {
		workSource = playbookSource
	}

	// Create container with Ansible
	container := m.buildContainer(rolesSource, "").
		WithDirectory("/ansible", rolesSource).
		WithDirectory("/work", workSource).
		WithWorkdir("/work/" + workdir).
		WithEnvVariable("ANSIBLE_ROLES_PATH", "/ansible/roles").
		WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
		WithMountedSecret("/root/.ssh/id_rsa", sshPrivateKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithExec([]string{"ssh-keyscan", "-H", "github.com"}, dagger.ContainerWithExecOpts{
			RedirectStdout: "/root/.ssh/known_hosts",
		})

	// Inject variables from WithVariable
	container, err := m.injectVariables(ctx, container)
	if err != nil {
		return "", err
	}

	// Run ansible-playbook
	container = container.WithExec(args)

	return container.Stdout(ctx)
}
