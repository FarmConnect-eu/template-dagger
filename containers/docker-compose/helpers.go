package main

import (
	"context"
	"dagger/docker-compose/internal/dagger"
	"fmt"
)

func (m *DockerCompose) buildContainer(
	ctx context.Context,
	source *dagger.Directory,
	composePath string,
) *dagger.Container {
	container := dag.Container().
		From("docker:27-cli").
		WithUnixSocket("/var/run/docker.sock", dag.Host().UnixSocket("/var/run/docker.sock")).
		WithMountedDirectory("/workspace", source).
		WithWorkdir("/workspace")

	container = container.WithExec([]string{
		"sh", "-c",
		"apk add --no-cache docker-cli-compose",
	})

	if m.RegistryHost != "" && m.RegistryUsername != nil && m.RegistryPassword != nil {
		container = container.
			WithSecretVariable("REGISTRY_USERNAME", m.RegistryUsername).
			WithSecretVariable("REGISTRY_PASSWORD", m.RegistryPassword).
			WithExec([]string{
				"sh", "-c",
				fmt.Sprintf("echo $REGISTRY_PASSWORD | docker login %s -u $REGISTRY_USERNAME --password-stdin", m.RegistryHost),
			})
	}

	for _, v := range m.Variables {
		if v.Secret != nil {
			container = container.WithSecretVariable(v.Key, v.Secret)
		} else {
			container = container.WithEnvVariable(v.Key, v.Value)
		}
	}

	return container
}

func getComposeCommand(composePath string) []string {
	if composePath == "" {
		composePath = "docker-compose.yml"
	}
	return []string{"docker", "compose", "-f", composePath}
}
