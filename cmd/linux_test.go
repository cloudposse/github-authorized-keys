package cmd

import (
	"testing"
	"strconv"
	"os/user"
	"fmt"
	"strings"
)

func inArray(val string, array []string) bool {
	for _, inArrayValue := range array {
		if val == inArrayValue {
			return true
		}
	}
	return false
}

func TestApiLinuxUserNotExists(t *testing.T) {
	t.Log("Check user not exists - Positive testing")

	userName := "test"

	isExists := linuxUserExists(userName)

	if isExists {
		t.Errorf("User should not exist.")
	}
}

func TestApiLinuxUserExists(t *testing.T) {
	t.Log("Check user exists - Positive testing")

	userName := "root"

	isExists := linuxUserExists(userName)

	if !isExists {
		t.Errorf("User should exist.")
	}
}

func TestApiLinuxCreateUser(t *testing.T) {
	t.Log("Create user - Positive testing")
	userName := linuxUser{Gid: "", Name: "test", Shell: "/bin/bash", Groups: []string{"wheel", "root"}}

	err := linuxUserCreate(userName)
	defer linuxUserDelete(userName)


	if err != nil {
		t.Errorf("User should be created, got error: %v", err)
	}


	osUser, _ := user.Lookup(userName.Name)

	if osUser.Username != userName.Name {
		t.Errorf("Linux user name %v should be equal %v", osUser.Username, userName.Name)
	}

	if value, _ := strconv.ParseInt(osUser.Gid, 10, 64); value <= 0 {
		t.Errorf("Linux user GID should be > 0, got %v", osUser.Gid)
	}

	gids, _:= osUser.GroupIds()

	for _, group := range userName.Groups {
		linuxGroup, err := user.LookupGroup(group)
		if err != nil {
			t.Errorf("Did not find group: %v. Got error %v", group, err)
		}


		if  ! inArray( string(linuxGroup.Gid), gids) {
			t.Errorf(fmt.Sprintf("Group %v does not contain user. User groups are %v", group,
				strings.Join(gids, ",")))
		}
	}

	shell := linuxUserShell(userName.Name)

	if ! strings.EqualFold(shell, userName.Shell) {
		t.Errorf("Expect user shell %v, got %v.", userName.Shell, shell)
	}
}

func TestApiLinuxCreateUserProvideGid(t *testing.T) {
	t.Log("Create user - Positive testing")
	userName := linuxUser{Gid: "42", Name: "test", Shell: "/bin/bash", Groups: []string{"root"}}

	err := linuxUserCreate(userName)
	defer linuxUserDelete(userName)


	if err != nil {
		t.Errorf("User should be created, got error: %v", err)
	}


	osUser, _ := user.Lookup(userName.Name)

	if osUser.Username != userName.Name {
		t.Errorf("Linux user name %v should be equal %v", osUser.Username, userName.Name)
	}

	if osUser.Gid != userName.Gid {
		t.Errorf("Linux user GID %v should be eqaul %v", osUser.Gid, userName.Gid)
	}

	gids, _:= osUser.GroupIds()

	for _, group := range userName.Groups {
		linuxGroup, err := user.LookupGroup(group)
		if err != nil {
			t.Errorf("Did not find group: %v. Got error %v", group, err)
		}


		if  ! inArray( string(linuxGroup.Gid), gids) {
			t.Errorf(fmt.Sprintf("Group %v does not contain user. User groups are %v", group,
				strings.Join(gids, ",")))
		}
	}

	shell := linuxUserShell(userName.Name)

	if ! strings.EqualFold(shell, userName.Shell) {
		t.Errorf("Expect user shell %v, got %v.", userName.Shell, shell)
	}
}


func TestApiLinuxGroupNotExists(t *testing.T) {
	t.Log("Check group not exists - Positive testing")

	groupName := "test"

	isExists := linuxGroupExists(groupName)

	if isExists {
		t.Errorf("Group should not exist.")
	}
}


func TestApiLinuxGroupExists(t *testing.T) {
	t.Log("Check group exists - Positive testing")

	groupName := "wheel"

	isExists := linuxGroupExists(groupName)

	if ! isExists {
		t.Errorf("Group should exist.")
	}
}


func TestApiLinuxGroupByIdNotExists(t *testing.T) {
	t.Log("Check group not exists - Positive testing")

	groupID := "43"

	isExists := linuxGroupExistsByID(groupID)

	if isExists {
		t.Errorf("Group should not exist.")
	}
}

func TestApiLinuxGroupByIdExists(t *testing.T) {
	t.Log("Check group exists - Positive testing")

	groupID := "42"

	isExists := linuxGroupExistsByID(groupID)

	if ! isExists {
		t.Errorf("Group should exist.")
	}
}


func TestApiLinuxUserShell(t *testing.T) {
	t.Log("Check getting user shell")

	userName := "root"

	shell := linuxUserShell(userName)

	if ! strings.EqualFold(shell, "/bin/ash") {
		t.Errorf("Expect user shell %v, got %v.", "/bin/ash", shell)
	}
}
