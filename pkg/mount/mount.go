package mount

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

const (
	hostDirectory = "/var/lib/runai/shared"
	containerDirectory = "/runai/shared"
)

func MountDirectoryToContainer(containerRootPath string) error {
	mountBinary, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	containerInternalPath := path.Join(containerRootPath, containerDirectory)
	args := []string{mountBinary, "--bind", hostDirectory, containerInternalPath}
	return syscall.Exec(mountBinary, args, os.Environ())
}
