// Package main provides a Dagger module for linting Dockerfiles with droast.
//
// droast (https://github.com/immanuwell/dockerfile-roast) is a standalone
// Dockerfile linter: no API key, no LLM, no network access at lint time.
// It complements security/trivy, which scans dependencies and secrets but does
// not review Dockerfile authoring practices.
//
// Severity model: droast exits 1 only on ERROR findings. WARNING/INFO findings
// are reported but leave the exit code at 0, so most projects get an advisory
// report rather than a gate.
package main

import (
	"context"
)

// defaultDroastVersion pins the linter image so rule changes land deliberately.
const defaultDroastVersion = "1.4.4"

type Droast struct {
	// Tag of the ghcr.io/immanuwell/droast image
	Version string
	// Lowest severity to report: info, warning, error (empty => droast default)
	MinSeverity string
	// Comma-separated rule IDs to ignore (e.g. "DF033,DF012")
	Skip string
	// Keep droast's joke lines in the output (off by default: noise in CI logs)
	Roast bool
}

func New() *Droast {
	return &Droast{
		Version: defaultDroastVersion,
	}
}

func (m *Droast) Test(ctx context.Context) (string, error) {
	return m.buildContainer().
		WithExec([]string{droastBin, "--version"}).
		Stdout(ctx)
}

// WithVersion overrides the droast image tag.
func (m *Droast) WithVersion(version string) *Droast {
	return &Droast{Version: version, MinSeverity: m.MinSeverity, Skip: m.Skip, Roast: m.Roast}
}

// WithMinSeverity filters out findings below the given severity: info, warning, error.
func (m *Droast) WithMinSeverity(severity string) *Droast {
	return &Droast{Version: m.Version, MinSeverity: severity, Skip: m.Skip, Roast: m.Roast}
}

// WithSkip ignores the given comma-separated rule IDs (e.g. "DF033,DF012").
func (m *Droast) WithSkip(rules string) *Droast {
	return &Droast{Version: m.Version, MinSeverity: m.MinSeverity, Skip: rules, Roast: m.Roast}
}

// WithRoast keeps droast's joke lines in the report.
func (m *Droast) WithRoast(
	// +optional
	// +default=true
	enabled bool,
) *Droast {
	return &Droast{Version: m.Version, MinSeverity: m.MinSeverity, Skip: m.Skip, Roast: enabled}
}
