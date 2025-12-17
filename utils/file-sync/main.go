// File Sync module for Dagger
//
// This module provides file synchronization to remote hosts via SSH.
// It allows copying local files and directories to remote destinations,
// replacing Ansible for simple configuration deployment scenarios.
//
// Example usage:
//
//	dagger call -m utils/file-sync \
//	  with-context --host 172.16.24.97 --user admin --ssh-key env:SSH_KEY \
//	  with-file --local-path settings/settings.yaml --remote-path /nfs/app/settings/settings.yaml \
//	  with-directory --local-path configs/ --remote-path /nfs/app/configs/ \
//	  sync --source .

package main

import (
	"dagger/file-sync/internal/dagger"
)

// FileSync module for synchronizing files to remote hosts via SSH
type FileSync struct {
	// SSH Context configuration
	SSHHost string
	SSHUser string
	SSHPort int
	SSHKey  *dagger.Secret

	// Files to synchronize
	Files []*FileMapping

	// Directories to synchronize
	Directories []*DirectoryMapping
}

// FileMapping represents a file to copy from local to remote
type FileMapping struct {
	LocalPath  string
	RemotePath string
	Owner      string // Optional owner in format "uid:gid" (e.g., "1001:1001")
	Mode       string // Optional file mode (e.g., "0644")
}

// DirectoryMapping represents a directory to copy from local to remote
type DirectoryMapping struct {
	LocalPath  string
	RemotePath string
	Owner      string // Optional owner in format "uid:gid" (e.g., "1001:1001")
	Mode       string // Optional directory mode (e.g., "0755")
	Recursive  bool   // Apply owner/mode recursively (default: true)
}

// New creates a new FileSync instance
func New() *FileSync {
	return &FileSync{
		SSHPort:     22,
		Files:       []*FileMapping{},
		Directories: []*DirectoryMapping{},
	}
}
