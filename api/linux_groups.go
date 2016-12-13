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
	"os/user"
)

const (
	// Group file contains one row per group
	// Format of the row consists at least of  3 columns
	// https://en.wikipedia.org/wiki/Group_identifier
	countOfColumnsInGroup = 3

	// Group name stored in 0 column
	nameColumnNumberInGroup = 0

	// Group password flag stored in 1 column
	passwordFlagColumnNumberInGroup = 1

	// Group GID stored in 2 column
	gidColumnNumberInGroup = 2

	// Group GID stored in 3 column
	usersColumnNumberInGroup = 3
)

func (linux *Linux) groupLookup(groupName string) (*user.Group, error) {
	groupInfo, err := linux.getEntity("group", groupName)

	if err != nil {
		return nil, user.UnknownGroupError(groupName)
	}

	if len(groupInfo) < countOfColumnsInGroup {
		return nil, errors.New("Wrong format of /etc/group")
	}

	group := user.Group{
		Gid:  groupInfo[gidColumnNumberInGroup],
		Name: groupInfo[nameColumnNumberInGroup],
	}

	return &group, err
}

func (linux *Linux) groupLookupByID(groupID string) (*user.Group, error) {
	groupInfo, err := linux.getEntity("group", groupID)

	if err != nil {
		return nil, user.UnknownGroupIdError(groupID)
	}

	if len(groupInfo) < countOfColumnsInGroup {
		return nil, errors.New("Wrong format of /etc/group")
	}

	group := user.Group{
		Gid:  groupInfo[gidColumnNumberInGroup],
		Name: groupInfo[nameColumnNumberInGroup],
	}

	return &group, err
}

// GroupExists - check if group {groupName} exists
func (linux *Linux) GroupExists(groupName string) bool {
	group, _ := linux.groupLookup(groupName)
	return group != nil
}

func (linux *Linux) groupExistsByID(groupID string) bool {
	group, _ := linux.groupLookupByID(groupID)
	return group != nil
}
