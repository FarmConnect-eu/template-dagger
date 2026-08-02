package main

// WithArg adds a Docker build argument (--build-arg)
//
// Build arguments are available during image build and can be referenced
// in Dockerfile using ARG instructions. Chain multiple calls for multiple args.
func (m *Docker) WithArg(
	// Argument name (e.g., "VERSION", "BUILD_DATE")
	key string,
	// Argument value
	value string,
) *Docker {
	c := m.clone()
	c.BuildArgs = append(c.BuildArgs, DockerBuildArg{Key: key, Value: value})
	return c
}
