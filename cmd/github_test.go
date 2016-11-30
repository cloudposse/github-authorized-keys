package cmd

import (
	"testing"
)

const (
	validToken = "7bf553ea09a665829455afd0f0541342fa85d71b"
	validOrg = "intervals-mining-lab"
	validTeamName = "libiada-developers"
	validTeamID = 191933
	validUser = "goruha"
)

func TestApiClient(t *testing.T) {
	t.Log("getTeam - Positive testing")

	token := validToken
	organization := validOrg
	teamName := validTeamName
	teamID := 0

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if team == nil {
		t.Errorf("Expected team object, got nil.")
	}


	if *team.ID == 0 {
		t.Errorf("Expected team id")
	}


	if *team.Name == "" {
		t.Errorf("Expected team name.")
	}
}

func TestApiClientGetTeamById(t *testing.T) {
	t.Log("getTeam - Positive testing get by ID")

	token := validToken
	organization := validOrg
	teamName := "dasdasd"
	teamID := validTeamID

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if team == nil {
		t.Errorf("Expected team object, got nil.")
	}


	if *team.ID == 0 {
		t.Errorf("Expected team id")
	}


	if *team.Name == "" {
		t.Errorf("Expected team name.")
	}
}


func TestApiClientWrongTeam(t *testing.T) {
	t.Log("getTeam - Wrong team testing")

	token := validToken
	organization := validOrg
	teamName := "xxx"
	teamID := 0

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	if err == nil {
		t.Errorf("Expected error, got %v.", err)
	}

	if err.Error() != "Team with such name or id not found" {
		t.Errorf("Wrong error, got %v.", err)
	}

	if team != nil {
		t.Errorf("Expected no team object, got %v.", team)
	}
}

func TestApiClientWrongToken(t *testing.T) {
	t.Log("getTeam - Wrong token testing")

	token := "11111111111111111111111111"
	organization := validOrg
	teamName := validTeamName
	teamID := 0

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	if err == nil {
		t.Errorf("Expected error, got %v.", err)
	}

	if err.Error() != "Access denied" {
		t.Errorf("Wrong error, got %v.", err)
	}

	if team != nil {
		t.Errorf("Expected no team object, got %v.", team)
	}
}

func TestApiClientWrongOrganization(t *testing.T) {
	t.Log("getTeam - Wrong organization testing")

	token := validToken
	organization := "dsadsad"
	teamName := validTeamName
	teamID := validTeamID

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	if err == nil {
		t.Errorf("Expected error, got %v.", err)
	}

	if err.Error() != "Access denied" {
		t.Errorf("Wrong error, got %v.", err)
	}

	if team != nil {
		t.Errorf("Expected no team object, got %v.", team)
	}
}

func TestApiClientIsMember(t *testing.T) {
	t.Log("IsMember - Positive testing")

	token := validToken
	organization := validOrg
	teamName := validTeamName
	teamID := validTeamID
	user := validUser

	c := newGithubClient(token, organization)
	team, _ := c.getTeam(teamName, teamID)

	isMember, err := c.isTeamMember(user, team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if !isMember {
		t.Errorf("User %v is member of team %v, but it was %v instead.", user, team, isMember)
	}
}

func TestApiClientIsNotMember(t *testing.T) {
	t.Log("IsMember - Not member")

	token := validToken
	organization := validOrg
	teamName := validTeamName
	teamID := validTeamID
	user := "dasda"

	c := newGithubClient(token, organization)
	team, _ := c.getTeam(teamName, teamID)

	isMember, err := c.isTeamMember(user, team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if isMember {
		t.Errorf("User %v is member of team %v, but it was %v instead.", user, team, isMember)
	}
}

func TestApiClientGetUser(t *testing.T) {
	t.Log("GetUser - Positive testing")

	token := validToken
	organization := validOrg
	userName := validUser

	c := newGithubClient(token, organization)
	user, err := c.getUser(userName)


	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if user == nil {
		t.Errorf("Expect to get user")
	}
}

func TestApiClientGetUserWrongUser(t *testing.T) {
	t.Log("GetUser - Wrong User")

	token := validToken
	organization := validOrg
	userName := "dasdddds232dasdas"

	c := newGithubClient(token, organization)
	user, err := c.getUser(userName)


	if err == nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if user != nil {
		t.Errorf("Expect to get nil, got %v", user)
	}
}


func TestApiClientGetPublicKeys(t *testing.T) {
	t.Log("GetPublicKeys - Positive testing")

	token := validToken
	organization := validOrg
	userName := validUser

	c := newGithubClient(token, organization)
	user, _ := c.getUser(userName)

	keys, err := c.getKeys(user)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if len(keys) <= 0 {
		t.Errorf("Expect to get keys, got %v.", keys)
	}
}

func TestApiClientGetTeamMembers(t *testing.T) {
	t.Log("GetTeamMembers - Positive testing")

	token := validToken
	organization := validOrg
	teamName := validTeamName
	teamID := 0

	c := newGithubClient(token, organization)
	team, err := c.getTeam(teamName, teamID)

	members, err := c.getTeamMembers(team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if len(members) <= 0 {
		t.Errorf("Expect to get users, got %v.", members)
	}
}


