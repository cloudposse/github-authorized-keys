package key_storages

import (
	"errors"
	"strings"
	log "github.com/Sirupsen/logrus"

	"github.com/cloudposse/github-authorized-keys/api"
)

type githubKeys struct {
	client *api.GithubClient
	team string
	teamID int
}

// Get {user} ssh keys
func (s *githubKeys) Get(user string) (value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			value = ""
			err = ErrStorageConnectionFailed
		}
	}()

	value = ""

	logger := log.WithFields(log.Fields{"class": "GithubClient", "method": "Get"})

	// Load team
	team, err := s.client.GetTeam(s.team, s.teamID)
	if err != nil {
		err = ErrStorageKeyNotFound
		return
	}

	// Check if user is a member
	isMember, err := s.client.IsTeamMember(user, team)
	if err != nil {
		err = ErrStorageKeyNotFound
		return
	}

	if ! isMember {
		err = ErrStorageKeyNotFound
		return
	}


	keys, response, err := s.client.GetKeys(user, nil)

	logger.Debugf("Response: %v", response)
	logger.Debugf("Response.StatusCode: %v", response.StatusCode)

	switch response.StatusCode {
		case 200:
			result := []string{}
			for _, value := range keys {
				result = append(result, *value.Key)
			}
			value = strings.Join(result, "\n")
			return

		case 404:
			value = ""
			err = ErrStorageKeyNotFound
			return

		default:
			value = ""
			err = errors.New("Access denied")
			return
	}
}

func NewGithubKeys(token, owner, team string, teamID int) *githubKeys {
	return &githubKeys{client: api.NewGithubClient(token, owner), team: team, teamID: teamID}
}