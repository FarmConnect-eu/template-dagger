package main

// WithDirectory adds a directory to synchronize from local source to remote destination
//
// Parameters:
//   - localPath: Path to the directory in the source (relative to --source)
//   - remotePath: Absolute path on the remote host where the directory contents should be placed
//   - owner: Optional owner in format "uid:gid" (e.g., "1001:1001")
//   - mode: Optional directory mode (e.g., "0755")
//   - recursive: Apply owner/mode recursively to all contents (default: true)
//
// The remote directory will be created automatically if it doesn't exist.
// All files and subdirectories will be copied recursively.
//
// Example:
//
//	dagger call with-context --host X --user Y --ssh-key env:KEY \
//	  with-directory --local-path configs/ --remote-path /app/configs/ --owner "1001:1001" \
//	  sync --source .
func (m *FileSync) WithDirectory(
	localPath string,
	remotePath string,
	// +optional
	owner string,
	// +optional
	mode string,
	// +optional
	// +default=true
	recursive bool,
) *FileSync {
	newDirs := append(copyDirectories(m.Directories), &DirectoryMapping{
		LocalPath:  localPath,
		RemotePath: remotePath,
		Owner:      owner,
		Mode:       mode,
		Recursive:  recursive,
	})

	return &FileSync{
		SSHHost:     m.SSHHost,
		SSHUser:     m.SSHUser,
		SSHPort:     m.SSHPort,
		SSHKey:      m.SSHKey,
		Files:       copyFiles(m.Files),
		Directories: newDirs,
	}
}
