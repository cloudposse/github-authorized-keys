package cmd

import (
	"os/user"
	"os/exec"
	"bytes"
	"io"
	"strings"
)

const(
	createUserCMD = "adduser"
	deleteUserCMD = "deluser"
)

type User struct {
	Name     string
	Gid      string // primary group ID
	Groups   []string
	Shell    string
}


func LinuxUserExists(user_name string) bool {
	user, _ := user.Lookup(user_name)
	return user != nil
}


func LinuxUserCreate(new User) error {
	var cmd *exec.Cmd

	if new.Gid == "" {
		cmd = exec.Command(createUserCMD, "-s", new.Shell, "-D", new.Name)
	} else {
		primaryGroup, err := user.LookupGroupId(new.Gid)
		if err != nil { return err }

		cmd = exec.Command(createUserCMD, "-s", new.Shell, "-G", primaryGroup.Name, "-D", new.Name)
	}

	err := cmd.Run()
	if err != nil { return err }


	for _, group := range new.Groups {
		cmd := exec.Command(createUserCMD, new.Name, group)
		err := cmd.Run()
		if err != nil { return err }
	}

	return nil
}

func LinuxUserDelete(new User) error {
	cmd := exec.Command(deleteUserCMD, new.Name)
	return cmd.Run()
}

func LinuxGroupExists(groupName string) bool {
	group, _ := user.LookupGroup(groupName)
	return group != nil
}

func LinuxGroupExistsById(groupId string) bool {
	group, _ := user.LookupGroupId(groupId)
	return group != nil
}

func LinuxUserShell(userName string) string {
	c1 := exec.Command("getent", "passwd", userName)
	// @TODO: Redo this golang code instead of cut command
	c2 := exec.Command("cut" , "-d:", "-f7")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()

	err := c1.Wait()
	if err != nil { return "" }

	w.Close()

	err = c2.Wait()

	if err != nil { return "" }

	// Command return \n in the end
	return strings.Replace(string(b2.Bytes()), "\n", "", 1)
}