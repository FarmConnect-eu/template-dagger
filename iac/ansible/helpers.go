package main

import (
	"context"
	"fmt"

	"dagger/ansible/internal/dagger"
)

func (m *Ansible) buildContainer(
	source *dagger.Directory,
) *dagger.Container {
	ansibleVersion := m.AnsibleVersion
	if ansibleVersion == "" {
		ansibleVersion = "2.17"
	}

	container := dag.Container().From("ubuntu:24.04")

	container = container.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{
			"apt-get", "install", "-y", "--no-install-recommends",
			"python3",
			"python3-pip",
			"python3-venv",
			"ssh",
			"curl",
			"ca-certificates",
			"git",
			"openssh-client",
			"sshpass",
		})

	container = container.
		WithExec([]string{"python3", "-m", "venv", "/opt/ansible-venv"}).
		WithNewFile("/tmp/requirements.txt", requirementsTxt).
		WithExec([]string{"/opt/ansible-venv/bin/pip", "install", "--no-cache-dir", "-r", "/tmp/requirements.txt"}).
		WithExec([]string{"/opt/ansible-venv/bin/pip", "install", "--no-cache-dir", fmt.Sprintf("ansible==%s", ansibleVersion)}).
		WithNewFile("/etc/ansible/ansible.cfg", ansibleCfg)

	if m.Requirements != nil {
		galaxyCmd := "/opt/ansible-venv/bin/ansible-galaxy install -r /tmp/requirements.yml"
		retryScript := galaxyInstallScript(galaxyCmd)
		container = container.
			WithMountedFile("/tmp/requirements.yml", m.Requirements).
			WithNewFile("/tmp/galaxy-install.sh", retryScript, dagger.ContainerWithNewFileOpts{Permissions: 0755}).
			WithExec([]string{"/bin/bash", "/tmp/galaxy-install.sh"})
	} else {
		communityCmd := "/opt/ansible-venv/bin/ansible-galaxy collection install community.general"
		posixCmd := "/opt/ansible-venv/bin/ansible-galaxy collection install ansible.posix"

		communityScript := galaxyInstallScript(communityCmd)
		posixScript := galaxyInstallScript(posixCmd)

		container = container.
			WithNewFile("/tmp/galaxy-install-community.sh", communityScript, dagger.ContainerWithNewFileOpts{Permissions: 0755}).
			WithExec([]string{"/bin/bash", "/tmp/galaxy-install-community.sh"}).
			WithNewFile("/tmp/galaxy-install-posix.sh", posixScript, dagger.ContainerWithNewFileOpts{Permissions: 0755}).
			WithExec([]string{"/bin/bash", "/tmp/galaxy-install-posix.sh"})
	}

	container = container.
		WithEnvVariable("PATH", "/opt/ansible-venv/bin:$PATH").
		WithEnvVariable("ANSIBLE_HOST_KEY_CHECKING", "False").
		WithEnvVariable("ANSIBLE_FORCE_COLOR", "true").
		WithEnvVariable("ANSIBLE_STDOUT_CALLBACK", "yaml").
		WithEnvVariable("PYTHONUNBUFFERED", "1")

	container = container.
		WithDirectory("/work", source).
		WithWorkdir("/work")

	container = container.
		WithExec([]string{"apt-get", "clean"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*", "/tmp/*"})

	return container
}

func (m *Ansible) injectVariables(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	for _, v := range m.Variables {
		if v.SecretValue != nil {
			container = container.WithSecretVariable(v.Key, v.SecretValue)
		} else {
			container = container.WithEnvVariable(v.Key, v.Value)
		}
	}

	return container, nil
}

func (m *Ansible) buildExtraVars() string {
	if len(m.ExtraVars) == 0 {
		return ""
	}

	vars := "{"
	first := true
	for _, kv := range m.ExtraVars {
		if !first {
			vars += ","
		}
		vars += fmt.Sprintf("\"%s\":\"%s\"", kv.Key, kv.Value)
		first = false
	}
	vars += "}"

	return vars
}

func galaxyInstallScript(command string) string {
	return fmt.Sprintf(`#!/bin/bash
set -e
MAX_RETRIES=3
RETRY_DELAY=2

for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $i of $MAX_RETRIES: %s"
    if %s; then
        echo "Success!"
        exit 0
    fi

    if [ $i -lt $MAX_RETRIES ]; then
        echo "Failed, retrying in ${RETRY_DELAY}s..."
        sleep $RETRY_DELAY
        RETRY_DELAY=$((RETRY_DELAY * 2))
    fi
done

echo "Failed after $MAX_RETRIES attempts"
exit 1
`, command, command)
}
