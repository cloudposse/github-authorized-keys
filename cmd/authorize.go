package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"github.com/spf13/viper"
)

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

		userName := args[0]

		c := newGithubClient(githubAPIToken, githubOrganization)

		// Load team
		team, err := c.getTeam(githubTeamName, githubTeamID)
		if err != nil { return err }

		// Check if user is a member
		isMember, err := c.isTeamMember(userName, team)
		if err != nil { return err }
		if ! isMember {
			return fmt.Errorf("User %v is not a member of team %v", userName, *team.Name)
		}

		// Load user
		user, err := c.getUser(userName)
		if err != nil { return err }
		// Get keys
		keys, err := c.getKeys(user)
		if err != nil { return err }

		// Print keys
		for _, k := range keys {
			fmt.Println(*k.Key)
		}

		return err
	},
}

func init() {
	RootCmd.AddCommand(authorizeCmd)
}
