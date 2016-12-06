package cmd

import (
	"fmt"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/cloudposse/github-authorized-keys/key_storages"
)

// Default ttl - 1day in seconds = 24 hours * 60 minutes * 60 seconds
const ETCDTTLDefault = int64(24 * 60 * 60)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize [user]",
	Short: "Outputs user public key if the user is member of github team",
	Long:
`
Outputs [user] public key if [user] is member of github team.

Could be used as provider for ssh AuthorizedKeysCommand.
To implement this add in /etc/ssh/sshd_config following string
-----------------------------------------------------------
|  AuthorizedKeysCommand github-authorized-keys authorize |
-----------------------------------------------------------
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		githubAPIToken 		:= viper.GetString("github_api_token")
		githubTeamName 		:= viper.GetString("github_team")
		githubTeamID 		:= viper.GetInt("github_team_id")
		githubOrganization 	:= viper.GetString("github_organization")


		// Validate user name arg
		if len(args) <= 0 {
			return errors.New("User name is required argument")
		}  else if len(args) > 1 {
			return errors.New("Can only authorize a single user at a time")
		}

		if githubAPIToken == "" {
			return errors.New("Github API Token is required")
		}

		if githubTeamName == "" && githubTeamID == 0 {
			return errors.New("Team name or Team id should be specified")
		}
		//-----------------------------------------------------------------


		var keys *key_storages.Proxy

		sourceStorage := key_storages.NewGithubKeys(githubAPIToken, githubOrganization, githubTeamName, githubTeamID)


		etcdEndpoints := []string{}
		if etcd := viper.GetString("etcd_endpoints"); etcd != "" {
			etcdEndpoints = strings.Split(etcd, ",")
		}

		if len(etcdEndpoints) > 0 {
			// add "s" suffix because duration should be in seconds.
			ttl, err := time.ParseDuration(viper.GetString("etcd_ttl") + "s")
			if err != nil {
				return fmt.Errorf("%v is not valid duration. %v", viper.GetString("etcd_ttl"), err)
			}

			fallbackStorage, _ := key_storages.NewEtcdCache(etcdEndpoints, ttl)
			keys = key_storages.NewProxy(sourceStorage, fallbackStorage)

		} else {

			fallbackStorage := &key_storages.NilStorage{}
			keys = key_storages.NewProxy(sourceStorage, fallbackStorage)
		}

		userName := args[0]


		// Get keys
		publicKeys, err := keys.Get(userName)
		if err != nil { return err }

		fmt.Println(publicKeys)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(authorizeCmd)

	authorizeCmd.Flags().StringSlice("etcd-endpoints", make([]string, 0),
		"Comma separeted gateways for etcd  ( environment variable ETCD could be used instead )")

	authorizeCmd.Flags().Int64("etcd-ttl", ETCDTTLDefault,
		"TTL sec for etcd cache ( environment variable ETCD_TTL could be used instead )")

	viper.BindPFlag("etcd_endpoints", authorizeCmd.Flags().Lookup("etcd-endpoints"))

	viper.BindPFlag("etcd_ttl", authorizeCmd.Flags().Lookup("etcd-ttl"))
}
