/*
 * Github Authorized Keys - Use GitHub teams to manage system user accounts and authorized_keys
 *
 * Copyright 2016 Cloud Posse, LLC <hello@cloudposse.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"bytes"
	"os/exec"
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
	if err != nil {
		return []string{}, err
	}

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
	logger.Debugf("Command:  %v", cmd)
	return linux.Command(cmd[0], cmd[1:]...)
}
