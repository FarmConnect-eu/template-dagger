package main

// WithTag adds image tags.
//
// If the tag is a semantic version (e.g., "v1.2.3", "v1.0.0-rc4"), it automatically
// generates all appropriate tags:
//   - v1.0.0-dev    -> v1.0.0-dev, dev
//   - v1.0.0-rc1    -> v1.0.0-rc1, v1.0.0-rc, v1.0-rc, v1-rc, rc
//   - v1.0.0-release -> v1.0.0-release, v1.0.0, v1.0, v1, release, latest
//   - v1.0.0        -> v1.0.0, v1.0, v1, latest
//
// For non-semver tags (e.g., "latest"), only that tag is added.
func (m *Docker) WithTag(
	// Image tag (e.g., "v1.2.3", "v1.0.0-rc4", "latest", "dev")
	tag string,
) *Docker {
	c := m.clone()
	c.Tags = append(c.Tags, parseSemanticVersion(tag)...)
	return c
}
