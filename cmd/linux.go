package cmd

import (
	"os/exec"
	"bytes"
	"strings"
)

const(
	chrootCommand     = "chroot"
)

type OS struct {
	root string
}

// Creates object allow interact with operating system
//
// rootDir - Path to directory contains linux root.
//
// Returns: OS object
func NewOs(rootDir string) OS {
	return OS{root: rootDir}
}

func (linux *OS) getEntity(database, key string) ([]string, error) {
	getent := linux.Command("getent", database, key)

	var b2 bytes.Buffer
	getent.Stdout = &b2

	err := getent.Run()
	if err != nil { return []string{}, err }

	row := strings.Trim(string(b2.Bytes()), "\n")

	columns := strings.Split(row, ":")
	return columns, nil
}

// Command returns the Cmd struct to execute the named program with
// the given arguments in context of OS
//
// It sets only the Path and Args in the returned structure.
//
// If name contains no path separators, Command uses LookPath to
// resolve the path to a complete name if possible. Otherwise it uses
// name directly.
//
// The returned Cmd's Args field is constructed from the command name
// followed by the elements of arg, so arg should not include the
// command name itself. For example, Command("echo", "hello")
func (linux *OS) Command(name string, arg ...string) *exec.Cmd {
	args := append( []string{linux.root, name}, arg...)
	return exec.Command(chrootCommand, args...)
}