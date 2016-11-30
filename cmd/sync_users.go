package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"strings"
	"github.com/spf13/viper"
)

// sync_usersCmd represents the sync_users command
var sync_usersCmd = &cobra.Command{
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


		githubApiToken		:= viper.GetString("github_api_token")
		githubTeamName 		:= viper.GetString("github_team")
		githubTeamId 		:= viper.GetInt("github_team_id")
		githubOrganization 	:= viper.GetString("github_organization")




		userGID 	:= viper.GetString("sync_users_gid")

		userGroups := []string{}
		if groups := viper.GetString("sync_users_groups"); groups != "" {
			userGroups = strings.Split(groups, ",")
		}

		userShell 	:= viper.GetString("sync_users_shell")

		if githubApiToken == "" {
			return errors.New("Github API Token is required")
		}

		// Validate Github Team exists
		if githubTeamName == "" && githubTeamId == 0 {
			return errors.New("Team name or Team id should be specified")
		}

		// If user GID is not empty validate that group with such id exists
		if userGID != "" && LinuxGroupExistsById(userGID) {
			return errors.New(fmt.Sprintf("Group with ID %v does not exists", userGID))
		}
		// Validate linux group exists
		nonExistedGroups := make([]string, 0)

		for _, group := range userGroups {
			if ! LinuxGroupExists(group) {
				nonExistedGroups = append(nonExistedGroups, group)
			}
		}

		if len(nonExistedGroups) > 0 {

			return errors.New(fmt.Sprintf("Groups %v not exists", strings.Join(nonExistedGroups, ",")))
		}

		//-------------------------------------------------------------------

		c := NewGithubClient(githubApiToken, githubOrganization)
		// Load team
		team, err := c.getTeam(githubTeamName, githubTeamId)
		if err != nil { return err }

		// Get all members
		githubUsers, err := c.GetTeamMembers(team)
		if err != nil { return err }

		// Here we will store user name for users that got error during creation
		notCreatedUsers := make([]string, 0)

		for _, githubUser := range githubUsers {
			// Create only non existed users
			if ! LinuxUserExists(*githubUser.Login) {

				linuxUser := User{Name: *githubUser.Login, Shell: userShell, Groups: userGroups}

				// If we have defined GID set it please
				if userGID != "" {
					linuxUser.Gid = userGID
				}

				// Create user and store it's name if there was error during creation
				if err := LinuxUserCreate(linuxUser); err != nil {
					notCreatedUsers = append(notCreatedUsers, linuxUser.Name)
				}
			}
 		}

		// Report error if we there was at least one error during user creation
		if len(notCreatedUsers) > 0 {
			return errors.New(fmt.Sprintf("Users %v created with errors", strings.Join(notCreatedUsers, ",")))
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(sync_usersCmd)

	sync_usersCmd.Flags().String("gid", "",
		"User's primary group id                       ( environment variable SYNC_USERS_GID    could be used instead )")

	sync_usersCmd.Flags().StringSlice("groups", make([]string, 0),
		"Comma separeted user's secondary groups name  ( environment variable SYNC_USERS_GROUPS could be used instead )")

	sync_usersCmd.Flags().String("shell", "/bin/bash",
		"User shell                                    ( environment variable SYNC_USERS_SHELL  could be used instead )")

	viper.BindPFlag("sync_users_gid", sync_usersCmd.Flags().Lookup("gid"))
	viper.BindPFlag("sync_users_groups",   sync_usersCmd.Flags().Lookup("groups"))
	viper.BindPFlag("sync_users_shell",  sync_usersCmd.Flags().Lookup("shell"))
}
