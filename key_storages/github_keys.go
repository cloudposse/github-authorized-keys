package keyStorages

import (
	"errors"
	"strings"

	"github.com/cloudposse/github-authorized-keys/api"
)

// GithubKeys - GitHub API as key storage
type GithubKeys struct {
	client *api.GithubClient
	team   string
	teamID int64
}

// Get - fetch {user} ssh keys
func (s *GithubKeys) Get(user string) (value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			value = ""
			err = ErrStorageConnectionFailed
		}
	}()

	value = ""

	// Load team
	team, err := s.client.GetTeam(s.team, s.teamID)
	if err != nil {
		if err == api.ErrorGitHubConnectionFailed {
			err = ErrStorageConnectionFailed
		} else {
			err = ErrStorageKeyNotFound
		}
		return
	}

	// Check if user is a member
	isMember, err := s.client.IsTeamMember(s.client.GetOrg(), user, team)
	if err != nil {
		if err == api.ErrorGitHubConnectionFailed {
			err = ErrStorageConnectionFailed
		} else {
			err = ErrStorageKeyNotFound
		}
		return
	}

	if !isMember {
		err = ErrStorageKeyNotFound
		return
	}

	keys, err := s.client.GetKeys(user)

	if err == nil {

		result := []string{}
		for _, value := range keys {
			result = append(result, *value.Key)
		}
		value = strings.Join(result, "\n")

	} else if err == api.ErrorGitHubNotFound {
		err = ErrStorageKeyNotFound
	} else {
		err = errors.New("access denied")
	}

	return
}

// NewGithubKeys - constructor for GitHub key storage
func NewGithubKeys(token, owner, githubURL, team string, teamID int64) *GithubKeys {
	return &GithubKeys{client: api.NewGithubClient(token, owner, githubURL), team: team, teamID: teamID}
}
