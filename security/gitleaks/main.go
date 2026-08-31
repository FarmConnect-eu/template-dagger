// Package main provides a Dagger module for secret scanning with gitleaks.
//
// It replaces the "curl the latest release then run the binary" step that CI
// workflows tend to inline: the version is pinned, the toolchain is a container
// image, and the same scan runs identically on a laptop and on a runner.
package main

import (
	"context"

	"dagger/gitleaks/internal/dagger"
)

// defaultGitleaksVersion is pinned on purpose: resolving `releases/latest` at
// job time makes CI break on a release nobody asked for, and adds an
// unauthenticated GitHub API call that is rate limited on shared runners.
const defaultGitleaksVersion = "8.30.1"

// findingsExitCode overrides the gitleaks default of 1. Gitleaks exits 1 both
// when it finds secrets and when it fails outright (bad config, unreadable
// repository), so the default value cannot tell a verdict from a crash. Moving
// findings to 2 leaves 1 meaning "gitleaks itself failed".
const findingsExitCode = 2

type Gitleaks struct {
	// Tag of the zricethezav/gitleaks image
	Version string
	// Custom gitleaks.toml replacing the built-in ruleset
	Config *dagger.File
	// Baseline of already-known findings to ignore
	Baseline *dagger.File
}

func New() *Gitleaks {
	return &Gitleaks{Version: defaultGitleaksVersion}
}

func (m *Gitleaks) Test(ctx context.Context) (string, error) {
	return m.buildContainer().
		WithExec([]string{"gitleaks", "version"}).
		Stdout(ctx)
}
