package cmd

import (
	"os/user"
	"os/exec"
)

const(
	createUserCMD = "adduser"
	updateUserCMD = "adduser"
	deleteUserCMD = "deluser"
)

func LinuxUserExists(user_name string) bool {
	user, _ := user.Lookup(user_name)
	return user != nil
}


func LinuxUserCreate(new user.User) error {
	cmd := exec.Command(createUserCMD, "-D", new.Name)
	return cmd.Run()
}

func LinuxUserDelete(new user.User) error {
	cmd := exec.Command(deleteUserCMD, new.Name)
	return cmd.Run()
}