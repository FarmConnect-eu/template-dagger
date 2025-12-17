package main

import (
	"dagger/docker-compose/internal/dagger"
)

// WithSecret adds a secret environment variable to inject into Docker Compose
//
// Secrets are handled securely by Dagger and never exposed in logs.
// Use this for sensitive values like API keys, passwords, database URIs.
//
// Parameters:
//   - key: Variable name (e.g., "MONGO_URI", "OPENAI_API_KEY")
//   - value: Secret value
//
// Example:
//
//	dagger call with-secret \
//	  --key MONGO_URI \
//	  --value env:MONGO_URI
func (m *DockerCompose) WithSecret(
	key string,
	value *dagger.Secret,
) *DockerCompose {
	newVars := copyVariables(m.Variables)
	newVars = append(newVars, &Variable{
		Key:    key,
		Value:  "",
		Secret: value,
	})

	return &DockerCompose{
		RegistryHost:     m.RegistryHost,
		RegistryUsername: m.RegistryUsername,
		RegistryPassword: m.RegistryPassword,
		Variables:        newVars,
		SSHHost:          m.SSHHost,
		SSHUser:          m.SSHUser,
		SSHPort:          m.SSHPort,
		SSHKey:           m.SSHKey,
		EnvFile:          m.EnvFile,
	}
}
