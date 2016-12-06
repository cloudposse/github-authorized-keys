package api

import (
	"errors"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Naive oauth setup

type accessToken struct {
	token *oauth2.Token
}

func (a accessToken) Token() (*oauth2.Token, error) {
	return a.token, nil
}

func newAccessToken(token string) accessToken {
	t := oauth2.Token{AccessToken: token}
	return accessToken{token: &t}
}

// GithubClient - client for operate with Github API
type GithubClient struct {
	client *github.Client
	owner  string
}

func (c *GithubClient) GetTeam(name string, id int) (*github.Team, error) {
	teams, response, err := c.client.Organizations.ListTeams(c.owner, nil)

	if response.StatusCode != 200 {
		return nil, errors.New("Access denied")
	}

	if err == nil {
		for _, team := range teams {
			if *team.ID == id || *team.Name == name {
				return team, err
			}
		}
	}
	return nil, errors.New("Team with such name or id not found")
}

func (c *GithubClient) getUser(name string) (*github.User, error) {
	user, response, err := c.client.Users.Get(name)

	if response.StatusCode != 200 {
		return nil, errors.New("Access denied")
	}

	return user, err
}

func (c *GithubClient) IsTeamMember(user string, team *github.Team) (bool, error) {
	result, _, err := c.client.Organizations.IsTeamMember(*team.ID, user)
	return result, err
}

//
func (c *GithubClient) GetKeys(userName string) ([]*github.Key, *github.Response, error) {
	return c.client.Users.ListKeys(userName, nil)
}


func (c *GithubClient) GetTeamMembers(team *github.Team) ([]*github.User, error) {
	users, _, err := c.client.Organizations.ListTeamMembers(*team.ID, nil)
	return users, err
}

func NewGithubClient(token, owner string) *GithubClient {
	c := oauth2.NewClient(oauth2.NoContext, newAccessToken(token))
	return &GithubClient{client: github.NewClient(c), owner: owner}
}