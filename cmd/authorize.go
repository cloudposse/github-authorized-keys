// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authorizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authorizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
