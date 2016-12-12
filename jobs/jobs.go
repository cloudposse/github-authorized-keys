package jobs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cloudposse/github-authorized-keys/api"
	"github.com/cloudposse/github-authorized-keys/config"
	"github.com/goruha/permbits"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"github.com/valyala/fasttemplate"
	"strings"
)

const wrapperScriptTpl = `#!/bin/bash
curl http://localhost:{port}/user/$1/authorized_keys
`

func init() {
	viper.SetDefault("ssh_restart_tpl", "/usr/sbin/service ssh force-reload")
}

func Run(cfg config.Config) {
	log.Info("Run syncUsers job on start")
	syncUsers(cfg)

	if cfg.IntegrateWithSSH {
		log.Info("Run ssh integration job on start")
		sshIntegrate(cfg)
	}

	if cfg.Interval != 0 {
		gocron.Every(cfg.Interval).Seconds().Do(syncUsers, cfg)

		// function Start start all the pending jobs
		gocron.Start()
		log.Info("Start jobs scheduler")
	}
}

func syncUsers(cfg config.Config) {
	logger := log.WithFields(log.Fields{"subsystem": "jobs", "job": "syncUsers"})

	linux := api.NewLinux(cfg.Root)

	c := api.NewGithubClient(cfg.GithubAPIToken, cfg.GithubOrganization)
	// Load team
	team, err := c.GetTeam(cfg.GithubTeamName, cfg.GithubTeamID)
	if err != nil {
		logger.Error(err)
		return
	}

	// Get all members
	githubUsers, err := c.GetTeamMembers(team)
	if err != nil {
		logger.Error(err)
		return
	}

	// Here we will store user name for users that got error during creation
	notCreatedUsers := make([]string, 0)

	for _, githubUser := range githubUsers {
		// Create only non existed users
		if !linux.UserExists(*githubUser.Login) {

			linuxUser := api.LinuxUser{Name: *githubUser.Login, Shell: cfg.UserShell, Groups: cfg.UserGroups}

			// If we have defined GID set it please
			if cfg.UserGID != "" {
				linuxUser.Gid = cfg.UserGID
			}

			// Create user and store it's name if there was error during creation
			if err := linux.UserCreate(linuxUser); err != nil {
				logger.Error(err)
				notCreatedUsers = append(notCreatedUsers, linuxUser.Name)
			}
		} else {
			logger.Debugf("User %v exists - skip creation", *githubUser.Login)
		}
	}
}

func sshIntegrate(cfg config.Config) {
	logger := log.WithFields(log.Fields{"subsystem": "jobs", "job": "sshIntegrate"})
	linux := api.NewLinux(cfg.Root)

	logger.Info("Ensure file /bin/github-authorized-keys")
	// Split listen string by : and get the port
	port := strings.Split(cfg.Listen, ":")[1]

	wrapperScript := fasttemplate.New(wrapperScriptTpl, "{", "}").
		ExecuteString(map[string]interface{}{"port": port})

	linux.FileEnsure("/bin/github-authorized-keys", wrapperScript)

	// Should be executable
	logger.Info("Ensure exec mode for file /bin/github-authorized-keys")
	linux.FileModeSet("/bin/github-authorized-keys", permbits.PermissionBits(0711))

	logger.Info("Ensure AuthorizedKeysCommand line in sshd_config")
	linux.FileEnsureLineMatch("/etc/ssh/sshd_config", "^AuthorizedKeysCommand\\s.*$", "AuthorizedKeysCommand /bin/github-authorized-keys")

	logger.Info("Ensure AuthorizedKeysCommandUser line in sshd_config")
	linux.FileEnsureLineMatch("/etc/ssh/sshd_config", "^AuthorizedKeysCommandUser\\s.*$", "AuthorizedKeysCommandUser nobody")

	logger.Info("Restart ssh")
	output, err := linux.TemplateCommand(viper.GetString("ssh_restart_tpl"), map[string]interface{}{}).CombinedOutput()
	logger.Infof("Output: %v", string(output))
	if err != nil {
		logger.Errorf("Error: %v", err.Error())
	}
}
