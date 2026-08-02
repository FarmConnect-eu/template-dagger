package main

// WithLabel adds an OCI image label (image config Config.Labels).
//
// Labels are applied at build time, visible via `docker inspect` on the host
// running the image, and preserved when the image is published to the registry.
// Chain multiple calls for multiple labels (e.g. org.opencontainers.image.*).
func (m *Docker) WithLabel(
	// Label name (e.g., "org.opencontainers.image.version")
	name string,
	// Label value
	value string,
) *Docker {
	c := m.clone()
	c.Labels = append(c.Labels, DockerBuildArg{Key: name, Value: value})
	return c
}
