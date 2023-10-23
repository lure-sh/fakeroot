package loginshell

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

var (
	// ErrInvalidPasswd is returned when the passwd file is invalid and can't be parsed
	ErrInvalidPasswd = errors.New("loginshell: invalid passwd file")

	// ErrNoSuchUser is returned when a given user isn't in the passwd file
	ErrNoSuchUser = errors.New("loginshell: provided uid not found in passwd file")

	// ErrUnsupported is returned on unsupported platforms, such as Windows
	ErrUnsupported = errors.New("loginshell: unsupported platform")
)

// Get returns the login shell belonging to the provided uid by parsing the passwd file.
// If uid is less than zero, the current uid will be used instead.
func Get(uid int) (string, error) {
	if uid < 0 {
		uid = os.Getuid()
	}

	// os.Getuid returns -1 on unsupported platforms
	if uid == -1 {
		return "", ErrUnsupported
	}

	fl, err := os.Open("/etc/passwd")
	if errors.Is(err, fs.ErrNotExist) {
		return "", ErrUnsupported
	} else if err != nil {
		return "", err
	}
	defer fl.Close()

	s := bufio.NewScanner(fl)
	for s.Scan() {
		luid, shell, err := parsePasswdLine(s.Text())
		if err != nil {
			return "", err
		}
		if luid == uid {
			return shell, nil
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return "", ErrNoSuchUser
}

func parsePasswdLine(line string) (int, string, error) {
	sline := strings.Split(line, ":")
	if len(sline) < 7 {
		return 0, "", ErrInvalidPasswd
	}
	uid, err := strconv.Atoi(sline[2])
	if err != nil {
		return 0, "", err
	}
	return uid, sline[6], nil
}
