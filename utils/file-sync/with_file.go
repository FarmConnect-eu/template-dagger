package main

// WithFile adds a file to synchronize from local source to remote destination
//
// Parameters:
//   - localPath: Path to the file in the source directory
//   - remotePath: Absolute path on the remote host where the file should be placed
//   - owner: Optional owner in format "uid:gid" (e.g., "1001:1001")
//   - mode: Optional file mode (e.g., "0644")
//
// The remote directory will be created automatically if it doesn't exist.
//
// Example:
//
//	dagger call with-context --host X --user Y --ssh-key env:KEY \
//	  with-file --local-path settings.yaml --remote-path /app/settings.yaml --owner "1001:1001" \
//	  sync --source .
func (m *FileSync) WithFile(
	localPath string,
	remotePath string,
	// +optional
	owner string,
	// +optional
	mode string,
) *FileSync {
	newFiles := append(copyFiles(m.Files), &FileMapping{
		LocalPath:  localPath,
		RemotePath: remotePath,
		Owner:      owner,
		Mode:       mode,
	})

	return &FileSync{
		SSHHost:     m.SSHHost,
		SSHUser:     m.SSHUser,
		SSHPort:     m.SSHPort,
		SSHKey:      m.SSHKey,
		Files:       newFiles,
		Directories: copyDirectories(m.Directories),
	}
}
