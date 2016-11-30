package cmd

import (
	"fmt"
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncUsersCmd represents the sync_users command
var syncUsersCmd = &cobra.Command{
	Use:   "sync_users",
	Short: "Create linux users for github team members",
	Long:
`Create user for each of github team member.
Run on schedule following command to create user asap.
-------------------------------------
|  github-authorized-keys sync_users|
-------------------------------------
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate Github API token


		githubAPIToken 		:= viper.GetString("github_api_token")
		githubTeamName 		:= viper.GetString("github_team")
		githubTeamID 		:= viper.GetInt("github_team_id")
		githubOrganization 	:= viper.GetString("github_organization")




		userGID 	:= viper.GetString("sync_users_gid")

		userGroups := []string{}
		if groups := viper.GetString("sync_users_groups"); groups != "" {
			userGroups = strings.Split(groups, ",")
		}

		userShell 	:= viper.GetString("sync_users_shell")

		if githubAPIToken == "" {
			return errors.New("Github API Token is required")
		}

		// Validate Github Team exists
		if githubTeamName == "" && githubTeamID == 0 {
			return errors.New("Team name or Team id should be specified")
		}

		// If user GID is not empty validate that group with such id exists
		if userGID != "" && linuxGroupExistsByID(userGID) {
			return fmt.Errorf("Group with ID %v does not exists", userGID)
		}
		// Validate linux group exists
		nonExistedGroups := make([]string, 0)

		for _, group := range userGroups {
			if ! linuxGroupExists(group) {
				nonExistedGroups = append(nonExistedGroups, group)
			}
		}

		if len(nonExistedGroups) > 0 {

			return fmt.Errorf("Groups %v not exists", strings.Join(nonExistedGroups, ","))
		}

		//-------------------------------------------------------------------

		c := newGithubClient(githubAPIToken, githubOrganization)
		// Load team
		team, err := c.getTeam(githubTeamName, githubTeamID)
		if err != nil { return err }

		// Get all members
		githubUsers, err := c.getTeamMembers(team)
		if err != nil { return err }

		// Here we will store user name for users that got error during creation
		notCreatedUsers := make([]string, 0)

		for _, githubUser := range githubUsers {
			// Create only non existed users
			if ! linuxUserExists(*githubUser.Login) {

				linuxUser := linuxUser{Name: *githubUser.Login, Shell: userShell, Groups: userGroups}

				// If we have defined GID set it please
				if userGID != "" {
					linuxUser.Gid = userGID
				}

				// Create user and store it's name if there was error during creation
				if err := linuxUserCreate(linuxUser); err != nil {
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

	syncUsersCmd.Flags().String("gid", "",
		"User's primary group id                       ( environment variable SYNC_USERS_GID    could be used instead )")

	syncUsersCmd.Flags().StringSlice("groups", make([]string, 0),
		"Comma separeted user's secondary groups name  ( environment variable SYNC_USERS_GROUPS could be used instead )")

	syncUsersCmd.Flags().String("shell", "/bin/bash",
		"User shell                                    ( environment variable SYNC_USERS_SHELL  could be used instead )")

	viper.BindPFlag("sync_users_gid", syncUsersCmd.Flags().Lookup("gid"))
	viper.BindPFlag("sync_users_groups",   syncUsersCmd.Flags().Lookup("groups"))
	viper.BindPFlag("sync_users_shell",  syncUsersCmd.Flags().Lookup("shell"))
}
