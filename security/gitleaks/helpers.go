package main

import (
	"fmt"
	"path"
	"strconv"

	"dagger/gitleaks/internal/dagger"
)

const (
	configPath   = "/etc/gitleaks/gitleaks.toml"
	baselinePath = "/etc/gitleaks/baseline.json"
	scanPath     = "/scan"
	reportDir    = "/out"
)

// buildContainer prepares a gitleaks container at the configured version.
//
// The safe.directory override is required, not defensive: gitleaks shells out to
// git, and git refuses to read a repository whose owner uid differs from the
// caller's ("detected dubious ownership"), which is always the case for a host
// checkout mounted into the container. It is set through GIT_CONFIG_* env vars
// rather than a `git config` exec so it costs no extra layer.
func (m *Gitleaks) buildContainer() *dagger.Container {
	version := m.Version
	if version == "" {
		version = defaultGitleaksVersion
	}

	container := dag.Container().
		From(fmt.Sprintf("zricethezav/gitleaks:v%s", version)).
		WithEnvVariable("GIT_CONFIG_COUNT", "1").
		WithEnvVariable("GIT_CONFIG_KEY_0", "safe.directory").
		WithEnvVariable("GIT_CONFIG_VALUE_0", "*")

	if m.Config != nil {
		container = container.WithMountedFile(configPath, m.Config)
	}
	if m.Baseline != nil {
		container = container.WithMountedFile(baselinePath, m.Baseline)
	}
	return container
}

// commonFlags builds the flags shared by every scan.
func (m *Gitleaks) commonFlags(redact bool) []string {
	flags := []string{
		"--no-banner",
		"--verbose",
		"--exit-code", strconv.Itoa(findingsExitCode),
	}
	// Bare --redact means 100%. It takes an optional value, so it must never be
	// passed as `--redact 100` (cobra would read 100 as a positional argument).
	if redact {
		flags = append(flags, "--redact")
	}
	if m.Config != nil {
		flags = append(flags, "--config", configPath)
	}
	if m.Baseline != nil {
		flags = append(flags, "--baseline-path", baselinePath)
	}
	return flags
}

// workdir resolves the directory the scan runs from.
func workdir(subpath string) string {
	if subpath == "" {
		subpath = "."
	}
	return path.Join(scanPath, subpath)
}

// hasFindings maps a gitleaks exit code to a verdict, and reports a non-zero
// code that is neither "clean" nor "leaks found" as the tool failure it is —
// otherwise a broken config would read as a passing scan.
func hasFindings(code int) (bool, error) {
	switch code {
	case 0:
		return false, nil
	case findingsExitCode:
		return true, nil
	default:
		return false, fmt.Errorf("gitleaks exited with code %d", code)
	}
}

// reportExtension maps a gitleaks report format to a file extension.
func reportExtension(format string) string {
	switch format {
	case "junit":
		return "xml"
	case "template":
		return "txt"
	default:
		// json, csv, sarif
		return format
	}
}
