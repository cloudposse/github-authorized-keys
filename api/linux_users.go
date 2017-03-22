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
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os/exec"
	"os/user"
	"github.com/cloudposse/github-authorized-keys/model/linux"
)

const (
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

func init() {
	// We need --force-badname because github users could contains capital letters, what is not acceptable in some distributions
	// Really regexp to verify badname rely on environment var that set in profile.d so we rarely hit this errors.
	//
	// adduser wants user name be the head and flags the tail.
	viper.SetDefault("linux_user_add_tpl", "adduser {username} --disabled-password --force-badname --shell {shell}")
	viper.SetDefault("linux_user_add_with_gid_tpl", "adduser {username} --disabled-password --force-badname --shell {shell} --group {group}")
	viper.SetDefault("linux_user_add_to_group_tpl", "adduser {username} {group}")
	viper.SetDefault("linux_user_del_tpl", "deluser {username}")
}

// UserExists - check if user {userName} exists
func (linux *Linux) UserExists(userName string) bool {
	user, _ := linux.userLookup(userName)
	return user != nil
}

func (linux *Linux) userLookup(userName string) (*user.User, error) {
	userInfo, err := linux.getEntity("passwd", userName)

	if err != nil {
		return nil, user.UnknownUserError(userName)
	}

	if len(userInfo) != countOfColumnsInPasswd {
		return nil, errors.New("Wrong format of /etc/passwd")
	}

	user := user.User{
		Uid:      userInfo[uidColumnNumberInPasswd],
		Gid:      userInfo[gidColumnNumberInPasswd],
		Name:     userInfo[nameColumnNumberInPasswd],
		Username: userInfo[nameColumnNumberInPasswd],
		HomeDir:  userInfo[homeColumnNumberInPasswd],
	}

	return &user, err
}

// UserCreate - create user {new}
func (linux *Linux) UserCreate(new linux.User) error {

	createUserCommandTemplate := viper.GetString("linux_user_add_tpl")
	createUserWithGIDCommandTemplate := viper.GetString("linux_user_add_with_gid_tpl")
	addUserToGroupCommandTemplate := viper.GetString("linux_user_add_to_group_tpl")

	var cmd *exec.Cmd

	template := createUserCommandTemplate

	args := map[string]interface{}{
		"shell":    new.Shell(),
		"username": new.Name(),
	}

	if new.Gid() != "" {
		template = createUserWithGIDCommandTemplate
		args["gid"] = new.Gid()

		if primaryGroup, err := linux.groupLookupByID(new.Gid()); err == nil {
			args["group"] = primaryGroup.Name
		}
	}

	cmd = linux.TemplateCommand(template, args)
	// cmd.Run called inside CombinedOutput()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%v\n", string(out))
		return err
	}

	fmt.Printf("Created user %v\n", new.Name())

	for _, group := range new.Groups() {
		cmd := linux.TemplateCommand(addUserToGroupCommandTemplate,
			map[string]interface{}{"username": new.Name(), "group": group})
		err := cmd.Run()
		if err != nil {
			return err
		}
		fmt.Printf("Added user %v to group %v\n", new.Name(), group)
	}

	return nil
}

func (linux *Linux) userDelete(new linux.User) error {
	deleteUserCommandTemplate := viper.GetString("linux_user_del_tpl")

	fmt.Printf("Delete user %v\n", new.Name())
	cmd := linux.TemplateCommand(deleteUserCommandTemplate, map[string]interface{}{"username": new.Name()})
	return cmd.Run()
}

func (linux *Linux) userShell(userName string) string {
	userInfo, err := linux.getEntity("passwd", userName)

	if err != nil || len(userInfo) != countOfColumnsInPasswd {
		return ""
	}

	return userInfo[shellColumnNumberInPasswd]
}
