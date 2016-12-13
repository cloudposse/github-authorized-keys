package config

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation"
	"time"
)

// Config - structure to store global configuration
type Config struct {
	GithubAPIToken     string
	GithubOrganization string
	GithubTeamName     string
	GithubTeamID       int

	EtcdEndpoints []string
	EtcdTTL       time.Duration
	EtcdPrefix    string

	UserGID    string
	UserGroups []string
	UserShell  string
	Root       string
	Interval   uint64

	IntegrateWithSSH bool

	Listen string
}

// Validate - process validation of config values
func (c *Config) Validate() (err error) {
	err = validation.StructRules{}.
		Add("GithubAPIToken", validation.Required.Error("is required")).
		Add("GithubOrganization", validation.Required.Error("is required")).
	// performs validation
		Validate(c)

	if err != nil {
		return
	}

	// Validate Github Team exists
	if c.GithubTeamName == "" && c.GithubTeamID == 0 {
		err = errors.New("Team name or Team id should be specified")
	}
	return
}
