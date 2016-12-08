package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cloudposse/github-authorized-keys/api"
)

// syncUsersCmd represents the sync-users command
var syncUsersCmd = &cobra.Command{
	Use:   "sync-users",
	Short: "Create linux users for github team members",
	Long: `Create user for each of github team member.
Run on schedule following command to create user asap.
-------------------------------------
|  github-authorized-keys sync-users|
-------------------------------------
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubAPIToken := viper.GetString("github_api_token")
		githubTeamName := viper.GetString("github_team")
		githubTeamID := viper.GetInt("github_team_id")
		githubOrganization := viper.GetString("github_organization")

		userGID := viper.GetString("sync_users_gid")

		userGroups := []string{}
		if groups := viper.GetString("sync_users_groups"); groups != "" {
			userGroups = strings.Split(groups, ",")
		}

		userShell := viper.GetString("sync_users_shell")

		if githubAPIToken == "" {
			return errors.New("Github API Token is required")
		}

		// Validate Github Team exists
		if githubTeamName == "" && githubTeamID == 0 {
			return errors.New("Team name or Team id should be specified")
		}

		root := viper.GetString("sync_users_root")

		linux := api.NewLinux(root)

		// Validate linux group exists
		nonExistedGroups := make([]string, 0)

		for _, group := range userGroups {
			if !linux.GroupExists(group) {
				nonExistedGroups = append(nonExistedGroups, group)
			}
		}

		if len(nonExistedGroups) > 0 {

			return fmt.Errorf("Groups not found: %v", strings.Join(nonExistedGroups, ","))
		}

		//-------------------------------------------------------------------

		c := api.NewGithubClient(githubAPIToken, githubOrganization)
		// Load team
		team, err := c.GetTeam(githubTeamName, githubTeamID)
		if err != nil {
			return err
		}

		// Get all members
		githubUsers, err := c.GetTeamMembers(team)
		if err != nil {
			return err
		}

		// Here we will store user name for users that got error during creation
		notCreatedUsers := make([]string, 0)

		for _, githubUser := range githubUsers {
			// Create only non existed users
			if !linux.UserExists(*githubUser.Login) {

				linuxUser := api.LinuxUser{Name: *githubUser.Login, Shell: userShell, Groups: userGroups}

				// If we have defined GID set it please
				if userGID != "" {
					linuxUser.Gid = userGID
				}

				// Create user and store it's name if there was error during creation
				if err := linux.UserCreate(linuxUser); err != nil {
					// @TODO: Replace with logger
					fmt.Printf("%v\n", err)
					notCreatedUsers = append(notCreatedUsers, linuxUser.Name)
				}
			}
		}

		// Report error if we there was at least one error during user creation
		if len(notCreatedUsers) > 0 {
			return fmt.Errorf("Users %v created with errors", strings.Join(notCreatedUsers, ","))
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncUsersCmd)

}
