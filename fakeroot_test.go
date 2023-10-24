package fakeroot_test

import (
	"errors"
	"os/exec"
	"syscall"
	"testing"

	"lure.sh/fakeroot"
)

func TestCommand(t *testing.T) {
	cmd, err := fakeroot.Command("whoami")
	if err != nil {
		t.Errorf("Unexpected error while creating command: %q", err)
	}

	data, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Unexpected error while executing command: %q", err)
	}
	if sdata := string(data); sdata != "root\n" {
		t.Errorf("Expected %q, got %q", "root\n", sdata)
	}
}

func TestCommandUIDError(t *testing.T) {
	cmd := exec.Command("whoami")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
	}
	err := fakeroot.Apply(cmd)
	if !errors.Is(err, fakeroot.ErrRootUIDAlreadyMapped) {
		t.Errorf("Expected error %q, got %q", fakeroot.ErrRootUIDAlreadyMapped, err)
	}
}

func TestCommandGIDError(t *testing.T) {
	cmd := exec.Command("whoami")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
	}
	err := fakeroot.Apply(cmd)
	if !errors.Is(err, fakeroot.ErrRootGIDAlreadyMapped) {
		t.Errorf("Expected error %q, got %q", fakeroot.ErrRootGIDAlreadyMapped, err)
	}
}
