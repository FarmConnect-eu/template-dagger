package main

// WithVariable adds an environment variable to inject into Docker Compose
//
// Use this for non-sensitive configuration values.
// For secrets, use WithSecret instead.
//
// Parameters:
//   - key: Variable name (e.g., "IMAGE_TAG", "PORT")
//   - value: Variable value
//
// Example:
//
//	dagger call with-variable --key IMAGE_TAG --value v1.2.3
func (m *DockerCompose) WithVariable(
	key string,
	value string,
) *DockerCompose {
	newVars := copyVariables(m.Variables)
	newVars = append(newVars, &Variable{
		Key:    key,
		Value:  value,
		Secret: nil,
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
