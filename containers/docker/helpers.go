package main

import (
	"fmt"
	"regexp"
	"strings"
)

// validateImageName ensures image name follows Docker conventions
func validateImageName(imageName string) error {
	if imageName == "" {
		return fmt.Errorf("image name cannot be empty")
	}
	if strings.ToLower(imageName) != imageName {
		return fmt.Errorf("image name must be lowercase: %s", imageName)
	}
	if strings.Contains(imageName, " ") {
		return fmt.Errorf("image name cannot contain spaces: %s", imageName)
	}
	return nil
}

// buildFullReference constructs the complete image reference
func buildFullReference(registryHost, imageName, tag string) string {
	if registryHost == "" {
		return fmt.Sprintf("%s:%s", imageName, tag)
	}
	return fmt.Sprintf("%s/%s:%s", registryHost, imageName, tag)
}

// getDefaultTags returns configured tags or "latest" if none
func (m *Docker) getDefaultTags() []string {
	if len(m.Tags) == 0 {
		return []string{"latest"}
	}
	return m.Tags
}

// parseSemanticVersion parses a semantic version and returns all applicable tags.
// For non-semver tags, returns just the original tag.
//
// Tagging strategy:
//   - v1.0.0-dev   -> v1.0.0-dev, v1.0.0-dev, v1.0-dev, v1-dev, dev
//   - v1.0.0-dev.1 -> v1.0.0-dev.1, v1.0.0-dev, v1.0-dev, v1-dev, dev
//   - v1.0.0-rc1   -> v1.0.0-rc1, v1.0.0-rc, v1.0-rc, v1-rc, rc
//   - v1.0.0-release -> v1.0.0-release, v1.0.0, v1.0, v1, release, latest
//   - v1.0.0       -> v1.0.0, v1.0, v1, latest
func parseSemanticVersion(version string) []string {
	var tags []string

	// Always add the exact version
	tags = append(tags, version)

	// Regex to parse semantic version: v1.2.3 or v1.2.3-suffix
	re := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-(.+))?$`)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		// Not a valid semver, return just the original tag
		return tags
	}

	major := matches[1]
	minor := matches[2]
	patch := matches[3]
	prerelease := matches[4]

	// Determine the version prefix (with or without 'v')
	prefix := ""
	if strings.HasPrefix(version, "v") {
		prefix = "v"
	}

	// Build version components
	patchTag := prefix + major + "." + minor + "." + patch
	minorTag := prefix + major + "." + minor
	majorTag := prefix + major

	// Handle prerelease suffixes
	if prerelease != "" {
		prereleaseLower := strings.ToLower(prerelease)

		if strings.HasPrefix(prereleaseLower, "dev") {
			// Dev tags: v1.0.0-dev, v1.0-dev, v1-dev, dev
			tags = append(tags, patchTag+"-dev", minorTag+"-dev", majorTag+"-dev", "dev")
		} else if strings.HasPrefix(prereleaseLower, "rc") {
			// RC tags: v1.0.0-rc, v1.0-rc, v1-rc, rc
			tags = append(tags, patchTag+"-rc", minorTag+"-rc", majorTag+"-rc", "rc")
		} else if prereleaseLower == "release" {
			// Release tags: v1.0.0, v1.0, v1, release, latest
			tags = append(tags, patchTag, minorTag, majorTag, "release", "latest")
		}
	} else {
		// No prerelease suffix = stable release: v1.0, v1, latest
		tags = append(tags, minorTag, majorTag, "latest")
	}

	return tags
}
