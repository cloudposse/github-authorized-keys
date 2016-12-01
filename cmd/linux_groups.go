package cmd

import (
	"os/user"
	"errors"
)

const(
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

func (linux *OS) groupLookup(groupName string) (*user.Group, error) {
	groupInfo, err := linux.getEntity("group", groupName)

	if err != nil {
		return nil, user.UnknownGroupError(groupName)
	}

	if len(groupInfo) < countOfColumnsInGroup {
		return nil, errors.New("Wrong format of /etc/group")
	}

	group := user.Group{
		Gid: groupInfo[gidColumnNumberInGroup],
		Name: groupInfo[nameColumnNumberInGroup],
	}

	return &group, err
}

func (linux *OS) groupLookupByID(groupID string) (*user.Group, error) {
	groupInfo, err := linux.getEntity("group", groupID)

	if err != nil {
		return nil, user.UnknownGroupIdError(groupID)
	}

	if len(groupInfo) < countOfColumnsInGroup {
		return nil, errors.New("Wrong format of /etc/group")
	}

	group := user.Group{
		Gid: groupInfo[gidColumnNumberInGroup],
		Name: groupInfo[nameColumnNumberInGroup],
	}

	return &group, err
}

func (linux *OS) groupExists(groupName string) bool {
	group, _ := linux.groupLookup(groupName)
	return group != nil
}

func (linux *OS) groupExistsByID(groupID string) bool {
	group, _ := linux.groupLookupByID(groupID)
	return group != nil
}