package cmd

import (
	"testing"
)

const (
	validToken = "7bf553ea09a665829455afd0f0541342fa85d71b"
	validOrg = "intervals-mining-lab"
	validTeamName = "libiada-developers"
	validTeamId = 191933
	validUser = "goruha"
)

func TestApiClient(t *testing.T) {
	t.Log("getTeam - Positive testing")

	token := validToken
	organization := validOrg
	team_name := validTeamName
	team_id := 0

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

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
	team_name := "dasdasd"
	team_id := validTeamId

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

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
	team_name := "xxx"
	team_id := 0

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

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
	team_name := validTeamName
	team_id := 0

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

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
	team_name := validTeamName
	team_id := validTeamId

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

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
	team_name := validTeamName
	team_id := validTeamId
	user := validUser

	c := NewGithubClient(token, organization)
	team, _ := c.getTeam(team_name, team_id)

	isMember, err := c.isTeamMember(user, team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if !isMember {
		t.Errorf("User %v is member of team %v, but it was %d instead.", user, team, isMember)
	}
}

func TestApiClientIsNotMember(t *testing.T) {
	t.Log("IsMember - Not member")

	token := validToken
	organization := validOrg
	team_name := validTeamName
	team_id := validTeamId
	user := "dasda"

	c := NewGithubClient(token, organization)
	team, _ := c.getTeam(team_name, team_id)

	isMember, err := c.isTeamMember(user, team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if isMember {
		t.Errorf("User %v is member of team %v, but it was %d instead.", user, team, isMember)
	}
}

func TestApiClientGetUser(t *testing.T) {
	t.Log("GetUser - Positive testing")

	token := validToken
	organization := validOrg
	user_name := validUser

	c := NewGithubClient(token, organization)
	user, err := c.getUser(user_name)


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
	user_name := "dasdddds232dasdas"

	c := NewGithubClient(token, organization)
	user, err := c.getUser(user_name)


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
	user_name := validUser

	c := NewGithubClient(token, organization)
	user, _ := c.getUser(user_name)

	keys, err := c.GetKeys(user)

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
	team_name := validTeamName
	team_id := 0

	c := NewGithubClient(token, organization)
	team, err := c.getTeam(team_name, team_id)

	members, err := c.GetTeamMembers(team)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
	if len(members) <= 0 {
		t.Errorf("Expect to get users, got %v.", members)
	}
}


