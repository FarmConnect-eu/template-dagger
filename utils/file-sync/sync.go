package main

import (
	"context"
	"dagger/file-sync/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
)

// Sync executes the file synchronization to the remote host using rsync
//
// This operation will:
// 1. Create remote directories as needed
// 2. Sync all configured files to their remote destinations
// 3. Sync all configured directories recursively
// 4. Apply ownership (chown) and permissions (chmod) via rsync options
//
// Parameters:
//   - source: Directory containing the local files to sync
//   - compress: Enable compression during transfer (default: true)
//   - delete: Delete extraneous files from destination (default: false)
//   - dryRun: Show what would be transferred without actually doing it (default: false)
//
// Example:
//
//	dagger call with-context --host X --user Y --ssh-key env:KEY \
//	  with-file --local-path settings.yaml --remote-path /app/settings.yaml --owner "1001:1001" \
//	  sync --source . --compress --delete
func (m *FileSync) Sync(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default=true
	compress bool,
	// +optional
	// +default=false
	delete bool,
	// +optional
	// +default=false
	dryRun bool,
) (string, error) {
	if m.SSHHost == "" || m.SSHKey == nil {
		return "", fmt.Errorf("SSH context not configured. Use with-context first")
	}

	if len(m.Files) == 0 && len(m.Directories) == 0 {
		return "", fmt.Errorf("no files or directories configured. Use with-file or with-directory")
	}

	container := m.buildContainer(source)

	var output strings.Builder
	output.WriteString("=== File Sync Operation (rsync) ===\n")
	output.WriteString(fmt.Sprintf("Target: %s@%s\n", m.SSHUser, m.SSHHost))
	output.WriteString(fmt.Sprintf("Options: compress=%t, delete=%t, dry-run=%t\n\n", compress, delete, dryRun))

	// Build SSH command for rsync
	sshCmd := fmt.Sprintf("ssh -p %d -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null", m.SSHPort)

	// Process files
	for _, f := range m.Files {
		remoteDir := filepath.Dir(f.RemotePath)
		output.WriteString(fmt.Sprintf("📁 Ensuring directory: %s\n", remoteDir))

		// Create remote directory first
		mkdirCmd := fmt.Sprintf("ssh -p %d %s@%s 'mkdir -p %s'",
			m.SSHPort, m.SSHUser, m.SSHHost, remoteDir)
		container = container.WithExec([]string{"sh", "-c", mkdirCmd})

		output.WriteString(fmt.Sprintf("📄 Syncing: %s → %s\n", f.LocalPath, f.RemotePath))

		// Build rsync command for file
		rsyncArgs := []string{"-av", "--progress"}

		if compress {
			rsyncArgs = append(rsyncArgs, "-z")
		}
		if dryRun {
			rsyncArgs = append(rsyncArgs, "--dry-run")
		}
		if f.Owner != "" {
			rsyncArgs = append(rsyncArgs, fmt.Sprintf("--chown=%s", f.Owner))
			output.WriteString(fmt.Sprintf("   👤 Owner: %s\n", f.Owner))
		}
		if f.Mode != "" {
			rsyncArgs = append(rsyncArgs, fmt.Sprintf("--chmod=%s", f.Mode))
			output.WriteString(fmt.Sprintf("   🔒 Mode: %s\n", f.Mode))
		}

		rsyncArgs = append(rsyncArgs, "-e", sshCmd)
		rsyncArgs = append(rsyncArgs, fmt.Sprintf("/workspace/%s", f.LocalPath))
		rsyncArgs = append(rsyncArgs, fmt.Sprintf("%s@%s:%s", m.SSHUser, m.SSHHost, f.RemotePath))

		container = container.WithExec(append([]string{"rsync"}, rsyncArgs...))
	}

	// Process directories
	for _, d := range m.Directories {
		remotePath := strings.TrimSuffix(d.RemotePath, "/")
		localPath := strings.TrimSuffix(d.LocalPath, "/")

		output.WriteString(fmt.Sprintf("📁 Ensuring directory: %s\n", remotePath))

		// Create remote directory first
		mkdirCmd := fmt.Sprintf("ssh -p %d %s@%s 'mkdir -p %s'",
			m.SSHPort, m.SSHUser, m.SSHHost, remotePath)
		container = container.WithExec([]string{"sh", "-c", mkdirCmd})

		output.WriteString(fmt.Sprintf("📂 Syncing directory: %s/ → %s/\n", localPath, remotePath))

		// Build rsync command for directory
		rsyncArgs := []string{"-av", "--progress"}

		if compress {
			rsyncArgs = append(rsyncArgs, "-z")
		}
		if delete {
			rsyncArgs = append(rsyncArgs, "--delete")
			output.WriteString("   🗑️  Delete extraneous files enabled\n")
		}
		if dryRun {
			rsyncArgs = append(rsyncArgs, "--dry-run")
		}
		if d.Owner != "" {
			if d.Recursive {
				rsyncArgs = append(rsyncArgs, fmt.Sprintf("--chown=%s", d.Owner))
			}
			output.WriteString(fmt.Sprintf("   👤 Owner: %s (recursive: %t)\n", d.Owner, d.Recursive))
		}
		if d.Mode != "" {
			if d.Recursive {
				// Use chmod for directories and files
				rsyncArgs = append(rsyncArgs, fmt.Sprintf("--chmod=D%s,F%s", d.Mode, d.Mode))
			}
			output.WriteString(fmt.Sprintf("   🔒 Mode: %s (recursive: %t)\n", d.Mode, d.Recursive))
		}

		rsyncArgs = append(rsyncArgs, "-e", sshCmd)
		// Trailing slash on source means "contents of directory"
		rsyncArgs = append(rsyncArgs, fmt.Sprintf("/workspace/%s/", localPath))
		rsyncArgs = append(rsyncArgs, fmt.Sprintf("%s@%s:%s/", m.SSHUser, m.SSHHost, remotePath))

		container = container.WithExec(append([]string{"rsync"}, rsyncArgs...))
	}

	// Execute and get output
	result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("sync failed: %w", err)
	}

	output.WriteString("\n--- rsync output ---\n")
	output.WriteString(result)
	output.WriteString("\n✅ Sync completed successfully\n")

	return output.String(), nil
}

// List shows the configured files and directories without executing
func (m *FileSync) List(
	ctx context.Context,
) (string, error) {
	var output strings.Builder
	output.WriteString("=== Configured Sync Items ===\n\n")

	if m.SSHHost != "" {
		output.WriteString(fmt.Sprintf("Target Host: %s@%s:%d\n\n", m.SSHUser, m.SSHHost, m.SSHPort))
	} else {
		output.WriteString("Target Host: Not configured (use with-context)\n\n")
	}

	output.WriteString("Files:\n")
	if len(m.Files) == 0 {
		output.WriteString("  (none)\n")
	}
	for _, f := range m.Files {
		ownerInfo := ""
		if f.Owner != "" {
			ownerInfo = fmt.Sprintf(" [owner: %s]", f.Owner)
		}
		modeInfo := ""
		if f.Mode != "" {
			modeInfo = fmt.Sprintf(" [mode: %s]", f.Mode)
		}
		output.WriteString(fmt.Sprintf("  📄 %s → %s%s%s\n", f.LocalPath, f.RemotePath, ownerInfo, modeInfo))
	}

	output.WriteString("\nDirectories:\n")
	if len(m.Directories) == 0 {
		output.WriteString("  (none)\n")
	}
	for _, d := range m.Directories {
		ownerInfo := ""
		if d.Owner != "" {
			recursiveInfo := ""
			if d.Recursive {
				recursiveInfo = ", recursive"
			}
			ownerInfo = fmt.Sprintf(" [owner: %s%s]", d.Owner, recursiveInfo)
		}
		modeInfo := ""
		if d.Mode != "" {
			modeInfo = fmt.Sprintf(" [mode: %s]", d.Mode)
		}
		output.WriteString(fmt.Sprintf("  📂 %s → %s%s%s\n", d.LocalPath, d.RemotePath, ownerInfo, modeInfo))
	}

	return output.String(), nil
}
