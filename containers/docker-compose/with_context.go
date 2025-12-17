package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithContext configures SSH-based remote Docker context for deployment
//
// This enables deployment to remote Docker hosts via SSH instead of local socket.
// When configured, all Docker commands will be executed on the remote host.
//
// Parameters:
//   - host: Remote host IP address or hostname
//   - user: SSH username for authentication
//   - port: SSH port (default: 22)
//   - sshKey: SSH private key for authentication
//
// Example:
//
//	dagger call with-context \
//	  --host 172.16.24.97 \
//	  --user admincd24 \
//	  --ssh-key env:SSH_PRIVATE_KEY \
//	  deploy --source . --project-name myapp
func (m *DockerCompose) WithContext(
	host string,
	user string,
	// +optional
	// +default=22
	port int,
	sshKey *dagger.Secret,
) *DockerCompose {
	if port == 0 {
		port = 22
	}

	return &DockerCompose{
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
		Variables:        copyVariables(m.Variables),
		SSHHost:          host,
		SSHUser:          user,
		SSHPort:          port,
		SSHKey:           sshKey,
		EnvFile:          m.EnvFile,
	}
}

// copyVariables creates a deep copy of the variables slice
func copyVariables(src []*Variable) []*Variable {
	if src == nil {
		return []*Variable{}
	}
	dst := make([]*Variable, len(src))
	for i, v := range src {
		dst[i] = &Variable{
			Key:    v.Key,
			Value:  v.Value,
			Secret: v.Secret,
		}
	}
	return dst
}
