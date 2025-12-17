package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithEnvFile mounts an environment file for Docker Compose
//
// The file will be mounted as .env in the workspace directory,
// making it available to docker-compose during deployment.
//
// Parameters:
//   - envFile: The .env file to mount
//
// Example:
//
//	dagger call with-env-file --env-file .env.production \
//	  deploy --source . --project-name myapp
func (m *DockerCompose) WithEnvFile(
	envFile *dagger.File,
) *DockerCompose {
	return &DockerCompose{
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
		Variables:        copyVariables(m.Variables),
		SSHHost:          m.SSHHost,
		SSHUser:          m.SSHUser,
		SSHPort:          m.SSHPort,
		SSHKey:           m.SSHKey,
		EnvFile:          envFile,
	}
}
