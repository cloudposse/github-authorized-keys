package cmd

import (
	"os/user"
	"os/exec"
	"fmt"
	"errors"
)

const(
	createUserCommand = "adduser"
	deleteUserCommand = "deluser"
)

const(
	// Passwd file contains one row per user
	// Format of the row consists of 7 columns
	// https://en.wikipedia.org/wiki/Passwd#Password_file
	countOfColumnsInPasswd = 7

	// User name stored in 0 column
	nameColumnNumberInPasswd = 0

	// User password flag stored in 1 column
	passwordFlagColumnNumberInPasswd = 1

	// User UID stored in 2 column
	uidColumnNumberInPasswd = 2

	// User GID stored in 3 column
	gidColumnNumberInPasswd = 3

	// User Gecos data stored in 4 column
	dataColumnNumberInPasswd = 4

	// User home dir stored in 5 column
	homeColumnNumberInPasswd = 5

	// User shell stored in 6 column
	shellColumnNumberInPasswd = 6

)

type linuxUser struct {
	Name     string
	Gid      string // primary group ID
	Groups   []string
	Shell    string
}

func (linux *OS) userExists(userName string) bool {
	user, _ := linux.userLookup(userName)
	return user != nil
}

func (linux *OS) userLookup(userName string) (*user.User, error) {
	userInfo, err := linux.getEntity("passwd", userName)

	if err != nil {
		return nil, user.UnknownUserError(userName)
	}

	if len(userInfo) != countOfColumnsInPasswd {
		return nil, errors.New("Wrong format of /etc/passwd")
	}

	user := user.User{
		Uid: userInfo[uidColumnNumberInPasswd],
		Gid: userInfo[gidColumnNumberInPasswd],
		Name: userInfo[nameColumnNumberInPasswd],
		Username: userInfo[nameColumnNumberInPasswd],
		HomeDir: userInfo[homeColumnNumberInPasswd],
	}

	return &user, err
}


func (linux *OS) userCreate(new linuxUser) error {
	var cmd *exec.Cmd

	userOptions := []string{
		"--shell", new.Shell,
		"--disabled-password",
		"--gecos", "''",
		new.Name,
	}

	if new.Gid != "" {
		primaryGroup, err := linux.groupLookupByID(new.Gid)
		if err != nil { return err }

		userOptions = append([]string{"--gid", primaryGroup.Name}, userOptions...)
	}
	cmd = linux.Command(createUserCommand, userOptions...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", cmd)
		return err
	}
	fmt.Printf("Created user %v\n", new.Name)

	for _, group := range new.Groups {
		cmd := linux.Command(createUserCommand, new.Name, group)
		err := cmd.Run()
		if err != nil { return err }
		fmt.Printf("Added user %v to group %v\n", new.Name, group)
	}

	return nil
}

func (linux *OS) userDelete(new linuxUser) error {
	cmd := linux.Command(deleteUserCommand, new.Name)
	return cmd.Run()
}

func (linux *OS) userShell(userName string) string {
	userInfo, err := linux.getEntity("passwd", userName)

	if err != nil || len(userInfo) != countOfColumnsInPasswd {
		return ""
	}

	return userInfo[shellColumnNumberInPasswd]
}