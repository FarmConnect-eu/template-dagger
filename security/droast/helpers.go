package main

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"

	"dagger/droast/internal/dagger"
)

// droastBin is the absolute path of the linter inside the image. Dagger's
// WithExec does not use the image entrypoint by default, so the binary is
// always invoked explicitly.
const droastBin = "/usr/local/bin/droast"

// dockerfileGlobs covers the naming variants in use across the workspace:
// `Dockerfile`, `Dockerfile.dev`, `api.Dockerfile`.
var dockerfileGlobs = []string{"**/Dockerfile", "**/Dockerfile.*", "**/*.Dockerfile"}

// prunedDirs are never linted: vendored or generated trees that may carry
// third-party Dockerfiles unrelated to the project being scanned.
var prunedDirs = []string{"node_modules", ".next", ".git", "vendor", "dist", "build"}

// emptySarif is a valid SARIF 2.1.0 document with no results, emitted when no
// Dockerfile exists so the artifact shape stays identical for CI consumers.
const emptySarif = `{
  "$schema": "https://json.schemastore.org/sarif-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": { "driver": { "name": "droast", "rules": [] } },
      "results": []
    }
  ]
}
`

// emptyReport returns a valid empty document for the requested format, used
// when no Dockerfile exists so the artifact shape never varies.
func emptyReport(format string) string {
	if format == "sarif" {
		return emptySarif
	}
	return "[]\n"
}

// buildContainer prepares a droast container at the configured version.
func (m *Droast) buildContainer() *dagger.Container {
	version := m.Version
	if version == "" {
		version = defaultDroastVersion
	}
	return dag.Container().
		From(fmt.Sprintf("ghcr.io/immanuwell/droast:%s", version))
}

// lintArgs builds the droast invocation shared by every operation.
func (m *Droast) lintArgs() []string {
	args := []string{droastBin}
	if !m.Roast {
		args = append(args, "--no-roast")
	}
	if m.MinSeverity != "" {
		args = append(args, "--min-severity", m.MinSeverity)
	}
	if m.Skip != "" {
		args = append(args, "--skip", m.Skip)
	}
	return args
}

// discoverDockerfiles lists the Dockerfiles under subpath, pruning generated
// and vendored trees. Paths are returned relative to the container workdir.
//
// Discovery happens here rather than by handing droast a directory because
// droast exits 1 when it finds no Dockerfile — indistinguishable from an ERROR
// finding. Globbing first lets the caller skip cleanly instead.
func discoverDockerfiles(ctx context.Context, source *dagger.Directory, subpath string) ([]string, error) {
	dir := source
	if subpath != "" && subpath != "." {
		dir = source.Directory(subpath)
	}

	seen := map[string]bool{}
	for _, pattern := range dockerfileGlobs {
		matches, err := dir.Glob(ctx, pattern)
		if err != nil {
			return nil, fmt.Errorf("globbing %q: %w", pattern, err)
		}
		for _, match := range matches {
			if isPruned(match) {
				continue
			}
			seen[match] = true
		}
	}

	files := make([]string, 0, len(seen))
	for file := range seen {
		files = append(files, file)
	}
	sort.Strings(files)
	return files, nil
}

// isPruned reports whether a matched path sits inside a pruned directory.
func isPruned(match string) bool {
	for _, segment := range strings.Split(path.Clean(match), "/") {
		for _, pruned := range prunedDirs {
			if segment == pruned {
				return true
			}
		}
	}
	return false
}

// shellQuote makes a path safe to embed in an `sh -c` command line.
func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

// shellCommand renders args as a single shell word list.
func shellCommand(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = shellQuote(arg)
	}
	return strings.Join(quoted, " ")
}

// reportExtension maps a droast output format to a report file extension.
func reportExtension(format string) string {
	switch format {
	case "sarif":
		return "sarif"
	case "json", "github":
		return "json"
	default:
		return "txt"
	}
}
