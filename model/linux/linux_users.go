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

package linux

import "strings"

// User - struct that extends os/user struct with shell param
type User struct {
	name   string
	gid    string // primary group ID
	groups []string
	shell  string
}

// NewUser - creates new User
func NewUser(name, gid string, groups []string, shell string) User {
	return User{name: name, gid: gid, groups: groups, shell: shell}
}

// Name - return user's name
func (user *User) Name() string {
	return strings.ToLower(user.name)
}

// Gid - return user gid
func (user *User) Gid() string {
	return user.gid
}

// Groups - return user groups
func (user *User) Groups() []string {
	return user.groups
}

// Shell - return user shell
func (user *User) Shell() string {
	return user.shell
}
