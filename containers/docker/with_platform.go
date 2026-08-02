package main

import (
	"dagger/docker/internal/dagger"
)

// WithPlatform sets the target platform for the build
//
// Overrides the default platform (linux/amd64).
func (m *Docker) WithPlatform(
	// Target platform (e.g., "linux/amd64", "linux/arm64")
	platform dagger.Platform,
) *Docker {
	c := m.clone()
	c.Platform = platform
	return c
}
