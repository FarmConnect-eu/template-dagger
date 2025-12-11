package main

import (
	"fmt"
	"regexp"
	"strings"
)

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

func buildFullReference(registryHost, imageName, tag string) string {
	if registryHost == "" {
		return fmt.Sprintf("%s:%s", imageName, tag)
	}
	return fmt.Sprintf("%s/%s:%s", registryHost, imageName, tag)
}

func (m *Docker) getDefaultTags() []string {
	if len(m.Tags) == 0 {
		return []string{"latest"}
	}
	return m.Tags
}

// parseSemanticVersion parses a semantic version and returns all applicable tags.
func parseSemanticVersion(version string) []string {
	var tags []string
	tags = append(tags, version)

	re := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-(.+))?$`)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return tags
	}

	major := matches[1]
	minor := matches[2]
	patch := matches[3]
	prerelease := matches[4]

	prefix := ""
	if strings.HasPrefix(version, "v") {
		prefix = "v"
	}

	patchTag := prefix + major + "." + minor + "." + patch
	minorTag := prefix + major + "." + minor
	majorTag := prefix + major

	if patchTag != version {
		tags = append(tags, patchTag)
	}
	tags = append(tags, minorTag, majorTag)

	if prerelease != "" {
		prereleaseLower := strings.ToLower(prerelease)
		if strings.HasPrefix(prereleaseLower, "rc") {
			tags = append(tags, "rc")
		} else if prereleaseLower == "release" {
			tags = append(tags, "release", "latest")
		}
	} else {
		tags = append(tags, "latest")
	}

	return tags
}
