package main

import (
	"context"
	"fmt"

	"dagger/droast/internal/dagger"
)

// Lint lints the Dockerfiles found under source and returns the report as text.
//
// Returns an error when droast reports ERROR-severity findings and
// failOnFindings is true. When no Dockerfile exists, it reports that and
// succeeds — an absent Dockerfile is not a lint failure.
func (m *Droast) Lint(
	ctx context.Context,
	// Directory to lint
	source *dagger.Directory,
	// Relative subpath within source (default: ".")
	// +optional
	// +default="."
	subpath string,
	// Output format: compact, json, sarif, github
	// +optional
	// +default="compact"
	format string,
	// Fail the call when ERROR-severity findings are detected
	// +optional
	// +default=false
	failOnFindings bool,
) (string, error) {
	if subpath == "" {
		subpath = "."
	}
	if format == "" {
		format = "compact"
	}

	files, err := discoverDockerfiles(ctx, source, subpath)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "No Dockerfile found - lint skipped\n", nil
	}

	args := append(m.lintArgs(), "--format", format)
	if !failOnFindings {
		args = append(args, "--no-fail")
	}
	args = append(args, files...)

	return m.workdir(source, subpath).WithExec(args).Stdout(ctx)
}

// LintArtifact lints the Dockerfiles under source and returns a directory
// holding `report.<ext>` (machine-readable), `report.txt` (human-readable) and
// `findings` ("1" when ERROR-severity findings exist, "0" otherwise).
//
// The call never fails on findings, so reports are always exportable and the CI
// gate is decided by reading `findings`. This mirrors trivy's ScanFsArtifact so
// both scanners are consumed the same way in a pipeline.
//
// When no Dockerfile exists the artifact is still produced, with an empty report
// and findings "0", so downstream publish steps do not need a special case.
func (m *Droast) LintArtifact(
	ctx context.Context,
	// Directory to lint
	source *dagger.Directory,
	// Relative subpath within source (default: ".")
	// +optional
	// +default="."
	subpath string,
	// Machine-readable report format: sarif, json, github (default: sarif)
	// +optional
	// +default="sarif"
	format string,
) (*dagger.Directory, error) {
	if subpath == "" {
		subpath = "."
	}
	if format == "" {
		format = "sarif"
	}

	files, err := discoverDockerfiles(ctx, source, subpath)
	if err != nil {
		return nil, err
	}

	reportName := "/out/report." + reportExtension(format)

	if len(files) == 0 {
		return dag.Directory().
			WithNewFile("report.txt", "No Dockerfile found - lint skipped\n").
			WithNewFile("report."+reportExtension(format), emptyReport(format)).
			WithNewFile("findings", "0").
			WithNewFile("dockerfiles", ""), nil
	}

	base := m.workdir(source, subpath).WithExec([]string{"mkdir", "-p", "/out"})

	// Human-readable run doubles as the verdict: `sh -c` propagates droast's
	// exit code through the redirect, and --no-fail is deliberately omitted.
	textCmd := shellCommand(append(append(m.lintArgs(), "--format", "compact"), files...)) + " > /out/report.txt"
	executed := base.WithExec(
		[]string{"sh", "-c", textCmd},
		dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny},
	)

	code, err := executed.ExitCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("droast lint: %w", err)
	}
	findings := "0"
	if code != 0 {
		findings = "1"
	}

	// Machine-readable run: --no-fail keeps the redirect exit code clean.
	machineCmd := shellCommand(append(append(m.lintArgs(), "--format", format, "--no-fail"), files...)) + " > " + reportName
	out := executed.WithExec([]string{"sh", "-c", machineCmd})

	return out.
		WithNewFile("/out/findings", findings).
		WithNewFile("/out/dockerfiles", joinLines(files)).
		Directory("/out"), nil
}

// workdir mounts source and positions the container on subpath.
func (m *Droast) workdir(source *dagger.Directory, subpath string) *dagger.Container {
	workdir := "/src"
	if subpath != "" && subpath != "." {
		workdir = "/src/" + subpath
	}
	return m.buildContainer().
		WithMountedDirectory("/src", source).
		WithWorkdir(workdir)
}

// joinLines renders a list as newline-terminated text.
func joinLines(values []string) string {
	out := ""
	for _, value := range values {
		out += value + "\n"
	}
	return out
}
