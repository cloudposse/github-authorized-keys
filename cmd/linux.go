package cmd

import (
	"os/user"
	"os/exec"
	"bytes"
	"io"
	"strings"
	"fmt"
)

const(
	createUserCommand = "adduser"
	deleteUserCommand = "deluser"
)

type linuxUser struct {
	Name     string
	Gid      string // primary group ID
	Groups   []string
	Shell    string
}


func linuxUserExists(userName string) bool {
	user, _ := user.Lookup(userName)
	return user != nil
}


func linuxUserCreate(new linuxUser) error {
	var cmd *exec.Cmd

	if new.Gid == "" {
		cmd = exec.Command(createUserCommand, "-s", new.Shell, "-D", new.Name)
	} else {
		primaryGroup, err := user.LookupGroupId(new.Gid)
		if err != nil { return err }

		cmd = exec.Command(createUserCommand, "-s", new.Shell, "-G", primaryGroup.Name, "-D", new.Name)
	}

	err := cmd.Run()
	if err != nil { return err }
	fmt.Printf("Created user %v\n", new.Name)

	for _, group := range new.Groups {
		cmd := exec.Command(createUserCommand, new.Name, group)
		err := cmd.Run()
		if err != nil { return err }
		fmt.Printf("Added user %v to group %v\n", new.Name, group)
	}

	return nil
}

func linuxUserDelete(new linuxUser) error {
	cmd := exec.Command(deleteUserCommand, new.Name)
	return cmd.Run()
}

func linuxGroupExists(groupName string) bool {
	group, _ := user.LookupGroup(groupName)
	return group != nil
}

func linuxGroupExistsByID(groupID string) bool {
	group, _ := user.LookupGroupId(groupID)
	return group != nil
}

func linuxUserShell(userName string) string {
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