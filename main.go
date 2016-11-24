package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
)

const (
	DEFAULT_CONFIG = "/etc/sshauth/config.json"
)

type Config struct {
	Token  string
	Owner  string `json:",omitempty"`
	Team   string `json:",omitempty"`
	TeamID int    `json:"team_id,omitempty"`
}

func loadConfig(file string) Config {
	f, err := os.Open(file)
	exitIf(err)

	decoder := json.NewDecoder(f)
	config := Config{}
	err = decoder.Decode(&config)
	exitIf(err)

	return config
}

func exitIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var configFile = flag.String("config", DEFAULT_CONFIG, "path to a JSON config file")
	flag.Parse()

	config := loadConfig(*configFile)

	c := NewGithubClient(config.Token, config.Owner)

	var (
		users []*github.User
		err   error
	)
	if config.TeamID != 0 {
		users, err = c.GetTeamMembersByID(config.TeamID)
	} else if config.Team != "" {
		users, err = c.GetTeamMembers(config.Team)
	} else {
		err = errors.New("Either team_id or team must be specified in config.json, but both were empty")
	}
	exitIf(err)

	keys := c.GetTeamKeys(users)
	for _, k := range keys {
		fmt.Println(*k.Key)
	}
}