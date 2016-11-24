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

var userGID int
var userGroups []string
var userShell string

// sync_usersCmd represents the sync_users command
var sync_usersCmd = &cobra.Command{
	Use:   "sync_users",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if githubApiToken == "" {
			return errors.New("Github API Token is required")
		}

		if githubTeamName == "" && githubTeamId == 0 {
			return errors.New("Team name or Team id should be specified")
		}

		// TODO: Work your own magic here
		fmt.Println("sync_users called")

		c := NewGithubClient(githubApiToken, githubOrganization)
		// Load team
		team, err := c.getTeam(githubTeamName, githubTeamId)
		if err != nil { return err }

		c.GetTeamMembers(team)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(sync_usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sync_usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sync_usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
