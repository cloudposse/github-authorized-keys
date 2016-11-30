package cmd

import (
	"os/user"
	"os/exec"
	"bytes"
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
	getent := exec.Command("getent", "passwd", userName)

	var b2 bytes.Buffer
	getent.Stdout = &b2

	err := getent.Run()
	if err != nil { return "" }

	userPasswd := strings.Replace(string(b2.Bytes()), "\n", "", 1)

	userPasswdSlice := strings.Split(userPasswd,":")

	if len(userPasswdSlice) != 7 {
		return ""
	}

	return userPasswdSlice[6]
}