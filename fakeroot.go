package fakeroot

import (
	"errors"
	"os"
	"os/exec"
	"slices"
	"syscall"
)

var (
	// ErrRootUIDAlreadyMapped is returned when there's already a mapping for the root user in a command
	ErrRootUIDAlreadyMapped = errors.New("fakeroot: root user has already been mapped in this command")

	// ErrRootGIDAlreadyMapped is returned when there's already a mapping for the root group in a command
	ErrRootGIDAlreadyMapped = errors.New("fakeroot: root group has already been mapped in this command")
)

// Command returns a command that runs in a fakeroot environment
func Command(name string, arg ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	return cmd, Apply(cmd)
}

// Apply applies the options required to run in a fakeroot environment to
// a command. It returns an error if the root group or user already has a mapping
// registered in the command.
func Apply(cmd *exec.Cmd) error {
	uid := os.Getuid()

	// If the user is already root, there's no need for fakeroot
	if uid == 0 {
		return nil
	}

	// Ensure SysProcAttr isn't nil
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Create a new user namespace
	cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER

	// If the command already contains a mapping for the root user, return an error
	if slices.ContainsFunc(cmd.SysProcAttr.UidMappings, rootMap) {
		return ErrRootUIDAlreadyMapped
	}

	// If the command already contains a mapping for the root group, return an error
	if slices.ContainsFunc(cmd.SysProcAttr.GidMappings, rootMap) {
		return ErrRootGIDAlreadyMapped
	}

	cmd.SysProcAttr.UidMappings = append(cmd.SysProcAttr.UidMappings, syscall.SysProcIDMap{
		ContainerID: 0,
		HostID:      uid,
		Size:        1,
	})

	cmd.SysProcAttr.GidMappings = append(cmd.SysProcAttr.GidMappings, syscall.SysProcIDMap{
		ContainerID: 0,
		HostID:      uid,
		Size:        1,
	})

	return nil
}

func rootMap(m syscall.SysProcIDMap) bool {
	return m.ContainerID == 0
}
