package main

import (
	"dagger/file-sync/internal/dagger"
)

// WithContext configures SSH connection for remote file synchronization
//
// Parameters:
//   - host: Remote host IP address or hostname
//   - user: SSH username for authentication
//   - port: SSH port (default: 22)
//   - sshKey: SSH private key for authentication
//
// Example:
//
//	dagger call with-context \
//	  --host 172.16.24.97 \
//	  --user admincd24 \
//	  --ssh-key env:SSH_PRIVATE_KEY \
//	  sync --source .
func (m *FileSync) WithContext(
	host string,
	user string,
	// +optional
	// +default=22
	port int,
	sshKey *dagger.Secret,
) *FileSync {
	if port == 0 {
		port = 22
	}

	return &FileSync{
		SSHHost:     host,
		SSHUser:     user,
		SSHPort:     port,
		SSHKey:      sshKey,
		Files:       copyFiles(m.Files),
		Directories: copyDirectories(m.Directories),
	}
}

// copyFiles creates a deep copy of the files slice
func copyFiles(src []*FileMapping) []*FileMapping {
	if src == nil {
		return []*FileMapping{}
	}
	dst := make([]*FileMapping, len(src))
	for i, f := range src {
		dst[i] = &FileMapping{
			LocalPath:  f.LocalPath,
			RemotePath: f.RemotePath,
			Owner:      f.Owner,
			Mode:       f.Mode,
		}
	}
	return dst
}

// copyDirectories creates a deep copy of the directories slice
func copyDirectories(src []*DirectoryMapping) []*DirectoryMapping {
	if src == nil {
		return []*DirectoryMapping{}
	}
	dst := make([]*DirectoryMapping, len(src))
	for i, d := range src {
		dst[i] = &DirectoryMapping{
			LocalPath:  d.LocalPath,
			RemotePath: d.RemotePath,
			Owner:      d.Owner,
			Mode:       d.Mode,
			Recursive:  d.Recursive,
		}
	}
	return dst
}
