package api

import (
	"os/exec"
	"bytes"
	"strings"

	"github.com/valyala/fasttemplate"
	"syscall"
)

// Linux - linux os with root dir
type Linux struct {
	root string
}

// NewLinux - creates object allow interact with operating system
//
// rootDir - Path to directory contains linux root.
//
// Returns: OS object
func NewLinux(rootDir string) Linux {
	return Linux{root: rootDir}
}

func (linux *Linux) getEntity(database, key string) ([]string, error) {
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
func (linux *Linux) Command(name string, params ...string) *exec.Cmd {
	cmd := exec.Command(name, params...)
	if strings.Trim(linux.root, " ") != "/" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Chroot: linux.root,
		}
	}

	return cmd
}

// TemplateCommand - creates command based on template and args with placeholders.
func (linux *Linux) TemplateCommand(template string, args map[string]interface{}) *exec.Cmd {
	t := fasttemplate.New(template, "{", "}")
	cmd := strings.Split(t.ExecuteString(args), " ")
	return linux.Command(cmd[0], cmd[1:]...)
}