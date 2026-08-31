package main

import (
	"context"
	"fmt"

	"dagger/gitleaks/internal/dagger"
)

// ScanGit scans the full git history of a repository (`gitleaks git`).
//
// The source directory must carry its .git directory, and the checkout must be
// complete (actions/checkout with fetch-depth: 0). Gitleaks replays `git log -p`,
// so a shallow clone silently narrows the scan to the commits that were fetched
// instead of failing — the scan stays green while covering almost nothing.
func (m *Gitleaks) ScanGit(
	ctx context.Context,
	// Repository to scan, including its .git directory
	source *dagger.Directory,
	// Relative subpath inside source (default: ".")
	// +optional
	// +default="."
	subpath string,
	// Redact detected secrets from the output
	// +optional
	// +default=true
	redact bool,
	// Fail the call when secrets are detected
	// +optional
	// +default=true
	failOnFindings bool,
	// git log options restricting the history, e.g. "--since=2024-01-01"
	// +optional
	logOpts string,
) (string, error) {
	executed := m.buildContainer().
		WithMountedDirectory(scanPath, source).
		WithWorkdir(workdir(subpath)).
		WithExec(m.gitArgs(redact, logOpts, "", ""),
			dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})

	code, err := executed.ExitCode(ctx)
	if err != nil {
		return "", fmt.Errorf("gitleaks scan: %w", err)
	}

	// Gitleaks logs its findings through zerolog, on stderr — stdout stays empty
	// unless a report is explicitly written to "-".
	output, err := executed.Stderr(ctx)
	if err != nil {
		return "", fmt.Errorf("read gitleaks output: %w", err)
	}

	findings, err := hasFindings(code)
	if err != nil {
		return "", fmt.Errorf("%w\n%s", err, output)
	}
	if findings && failOnFindings {
		return "", fmt.Errorf("gitleaks detected secrets in the git history:\n%s", output)
	}
	return output, nil
}

// ScanGitArtifact scans the git history and returns a directory holding the
// report (`report.<ext>`), the human-readable log (`report.txt`) and a
// `findings` file worth "1" when secrets were detected, "0" otherwise.
//
// Like the Trivy *-artifact variants, it never fails on findings: the report is
// always produced and exportable, and the CI gate is decided by reading
// `findings`. A genuine gitleaks failure is still returned as an error.
func (m *Gitleaks) ScanGitArtifact(
	ctx context.Context,
	// Repository to scan, including its .git directory
	source *dagger.Directory,
	// Relative subpath inside source (default: ".")
	// +optional
	// +default="."
	subpath string,
	// Redact detected secrets from the report
	// +optional
	// +default=true
	redact bool,
	// Report format: json, csv, junit, sarif (default: sarif)
	// +optional
	// +default="sarif"
	format string,
	// git log options restricting the history, e.g. "--since=2024-01-01"
	// +optional
	logOpts string,
) (*dagger.Directory, error) {
	if format == "" {
		format = "sarif"
	}
	reportPath := reportDir + "/report." + reportExtension(format)

	executed := m.buildContainer().
		WithMountedDirectory(scanPath, source).
		// Gitleaks writes the report file but does not create its parent
		// directory. Seeding an empty one avoids an extra mkdir exec.
		WithDirectory(reportDir, dag.Directory()).
		WithWorkdir(workdir(subpath)).
		WithExec(m.gitArgs(redact, logOpts, format, reportPath),
			dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})

	code, err := executed.ExitCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("gitleaks scan: %w", err)
	}

	output, err := executed.Stderr(ctx)
	if err != nil {
		return nil, fmt.Errorf("read gitleaks output: %w", err)
	}

	findings, err := hasFindings(code)
	if err != nil {
		return nil, fmt.Errorf("%w\n%s", err, output)
	}

	verdict := "0"
	if findings {
		verdict = "1"
	}

	return executed.
		WithNewFile(reportDir+"/report.txt", output).
		WithNewFile(reportDir+"/findings", verdict).
		Directory(reportDir), nil
}

// gitArgs assembles the `gitleaks git` command line. An empty format skips the
// report flags, keeping the scan log-only.
func (m *Gitleaks) gitArgs(redact bool, logOpts, format, reportPath string) []string {
	args := append([]string{"gitleaks", "git"}, m.commonFlags(redact)...)
	if format != "" {
		args = append(args, "--report-format", format, "--report-path", reportPath)
	}
	if logOpts != "" {
		args = append(args, "--log-opts", logOpts)
	}
	// Scan the working directory; the workdir already points at the subpath.
	return append(args, ".")
}
