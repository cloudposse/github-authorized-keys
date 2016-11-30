package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string


// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "github-authorized-keys",
	Short: "Allow to provide ssh access to servers based on github teams",
	Long:
`
github-authorized-keys is CLI tool allow to provide ssh access to server based on github teams.

Config:
  REQUIRED: Github API token        | flag --token   OR environment variable GITHUB_API_TOKEN
  REQUIRED: Github organization     | flag --org     OR environment variable GITHUB_ORGANIZATION
  REQUIRED: One of
  		   Github team name | flag --team    OR environment variable GITHUB_TEAM
  			OR
  		   Github team id   | flag --team_id OR Environment variable GITHUB_TEAM_ID
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

		// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-authorized-keys.yaml)")
	RootCmd.PersistentFlags().String("token", "", "Github API token (read more https://github.com/blog/1509-personal-api-tokens)")
	RootCmd.PersistentFlags().String("org", "", "Github organization")
	RootCmd.PersistentFlags().String("team", "", "Github team name")
	RootCmd.PersistentFlags().Int("team_id", 0, "Github team id")

	viper.BindPFlag("github_api_token",   RootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("github_organization",     RootCmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("github_team",    RootCmd.PersistentFlags().Lookup("team"))
	viper.BindPFlag("github_team_id", RootCmd.PersistentFlags().Lookup("team_id"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".github-authorized-keys") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
