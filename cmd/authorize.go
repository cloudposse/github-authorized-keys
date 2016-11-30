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

		githubApiToken		:= viper.GetString("github_api_token")
		githubTeamName 		:= viper.GetString("github_team")
		githubTeamId 		:= viper.GetInt("github_team_id")
		githubOrganization 	:= viper.GetString("github_organization")

		// Validate user name arg
		if len(args) <= 0 {
			return errors.New("User name is required argument")
		}  else if len(args) > 1 {
			return errors.New(
				fmt.Sprintf("Command does not support multiple users or you provide wrong user: %v", args))
		}

		if githubApiToken == "" {
			return errors.New("Github API Token is required")
		}

		if githubTeamName == "" && githubTeamId == 0 {
			return errors.New("Team name or Team id should be specified")
		}
		//-----------------------------------------------------------------

		user_name := args[0]

		c := NewGithubClient(githubApiToken, githubOrganization)

		// Load team
		team, err := c.getTeam(githubTeamName, githubTeamId)
		if err != nil { return err }

		// Check if user is a member
		isMember, err := c.isTeamMember(user_name, team)
		if err != nil { return err }
		if ! isMember {
			return errors.New(fmt.Sprintf("User %v is not a member of team %v", user_name, *team.Name))
		}

		// Load user
		user, err := c.getUser(user_name)
		if err != nil { return err }
		// Get keys
		keys, err := c.GetKeys(user)
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
