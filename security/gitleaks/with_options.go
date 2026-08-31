package main

import "dagger/gitleaks/internal/dagger"

// WithVersion sets the gitleaks image tag (immutable pattern).
func (m *Gitleaks) WithVersion(version string) *Gitleaks {
	return &Gitleaks{
		Version:  version,
		Config:   m.Config,
		Baseline: m.Baseline,
	}
}

// WithConfig mounts a custom gitleaks.toml, replacing the built-in ruleset.
func (m *Gitleaks) WithConfig(
	// gitleaks.toml to use instead of the default ruleset
	config *dagger.File,
) *Gitleaks {
	return &Gitleaks{
		Version:  m.Version,
		Config:   config,
		Baseline: m.Baseline,
	}
}

// WithBaseline mounts a baseline report whose findings are ignored.
//
// Use it to adopt gitleaks on a repository with known historical findings
// without turning the job red on day one.
func (m *Gitleaks) WithBaseline(
	// JSON report produced by a previous scan
	baseline *dagger.File,
) *Gitleaks {
	return &Gitleaks{
		Version:  m.Version,
		Config:   m.Config,
		Baseline: baseline,
	}
}
