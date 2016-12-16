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

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	// ErrorGitHubConnectionFailed - returned when there was connection error github.com
	ErrorGitHubConnectionFailed = errors.New("Connection to github.com failed")

	// ErrorGitHubAccessDenied - returned when there was access denied to github.com resource
	ErrorGitHubAccessDenied = errors.New("Access denied")
)

// Naive oauth setup
func newAccessToken(token string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
}

// GithubClient - client for operate with Github API
type GithubClient struct {
	client *github.Client
	owner  string
}

// GetTeam - return team structure based on name or id
func (c *GithubClient) GetTeam(name string, id int) (team *github.Team, err error) {
	defer func() {
		if r := recover(); r != nil {
			team = nil
			err = ErrorGitHubConnectionFailed
		}
	}()

	team = nil
	err = nil

	teams, response, _ := c.client.Organizations.ListTeams(c.owner, nil)

	if response.StatusCode != 200 {
		err = ErrorGitHubAccessDenied

	} else {
		for _, localTeam := range teams {
			if *localTeam.ID == id || *localTeam.Name == name {
				team = localTeam
				// team found
				return
			}
		}
		err = errors.New("Team with such name or id not found")
	}
	// Exit with error
	return
}

func (c *GithubClient) getUser(name string) (*github.User, error) {
	user, response, err := c.client.Users.Get(name)

	if response.StatusCode != 200 {
		return nil, ErrorGitHubAccessDenied
	}

	return user, err
}

// IsTeamMember - check if {user} is a membmer of {team}
func (c *GithubClient) IsTeamMember(user string, team *github.Team) (bool, error) {
	result, _, err := c.client.Organizations.IsTeamMember(*team.ID, user)
	return result, err
}

// GetKeys - return array of user's {userName} public keys
func (c *GithubClient) GetKeys(userName string) ([]*github.Key, *github.Response, error) {
	return c.client.Users.ListKeys(userName, nil)
}

// GetTeamMembers - return array of user's that are {team} members
func (c *GithubClient) GetTeamMembers(team *github.Team) (users []*github.User, err error) {
	defer func() {
		if r := recover(); r != nil {
			users = make([]*github.User, 0)
			err = ErrorGitHubConnectionFailed
		}
	}()

	users, _, err = c.client.Organizations.ListTeamMembers(*team.ID, nil)
	return
}

// NewGithubClient - constructor of GithubClient structure
func NewGithubClient(token, owner string) *GithubClient {
	c := oauth2.NewClient(oauth2.NoContext, newAccessToken(token))
	return &GithubClient{client: github.NewClient(c), owner: owner}
}
