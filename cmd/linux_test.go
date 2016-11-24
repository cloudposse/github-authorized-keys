package cmd

import (
	"testing"
	"os/user"
)

func TestApiLinuxUserNotExists(t *testing.T) {
	t.Log("Check user not exists - Positive testing")

	user_name := "test"

	isExists := LinuxUserExists(user_name)

	if isExists {
		t.Errorf("User should not exist.")
	}
}

func TestApiLinuxUserExists(t *testing.T) {
	t.Log("Check user exists - Positive testing")

	user_name := "root"

	isExists := LinuxUserExists(user_name)

	if !isExists {
		t.Errorf("User should exist.")
	}
}

func TestApiLinuxCreateUser(t *testing.T) {
	t.Log("Create user - Positive testing")

	user_name := user.User{Uid: "", Gid: "", Username: "test", Name: "test", HomeDir: "" }

	err := LinuxUserCreate(user_name)
	defer LinuxUserDelete(user_name)


	if err != nil {
		t.Errorf("User should be created, got error: %v", err)
	}
}
