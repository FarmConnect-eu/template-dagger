package main

// WithTarget sets the build target for multi-stage Dockerfiles
//
// Specifies which stage to build when using multi-stage builds.
func (m *Docker) WithTarget(
	// Stage name from Dockerfile (e.g., "builder", "production")
	target string,
) *Docker {
	c := m.clone()
	c.Target = target
	return c
}
