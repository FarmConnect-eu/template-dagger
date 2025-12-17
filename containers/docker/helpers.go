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

	// Add version hierarchy tags (without prerelease suffix)
	patchTag := prefix + major + "." + minor + "." + patch
	minorTag := prefix + major + "." + minor
	majorTag := prefix + major

	// Only add patch tag if different from the exact version
	if patchTag != version {
		tags = append(tags, patchTag)
	}
	tags = append(tags, minorTag, majorTag)

	// Handle prerelease suffixes
	if prerelease != "" {
		prereleaseLower := strings.ToLower(prerelease)

		if strings.HasPrefix(prereleaseLower, "rc") {
			tags = append(tags, "rc")
		} else if prereleaseLower == "release" {
			tags = append(tags, "release", "latest")
		}
	} else {
		// No prerelease suffix = stable release, add "latest"
		tags = append(tags, "latest")
	}

	return tags
}
