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
	Short: "Use GitHub teams to manage system user accounts and authorized_keys",
	Long: `
Use GitHub teams to manage system user accounts and authorized_keys.

Config:
  REQUIRED: Github API token        | flag --github-api-token    OR environment variable GITHUB_API_TOKEN
  REQUIRED: Github organization     | flag --github-organization OR environment variable GITHUB_ORGANIZATION
  REQUIRED: One of
  		   Github team name | flag --github-team    OR environment variable GITHUB_TEAM
  			OR
  		   Github team id   | flag --github-team-id OR Environment variable GITHUB_TEAM_ID
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

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-authorized-keys.yaml)")
	RootCmd.PersistentFlags().StringP("github-api-token", "t", "", "Github API token (read more https://github.com/blog/1509-personal-api-tokens)")
	RootCmd.PersistentFlags().StringP("github-organization", "o", "", "Github organization")
	RootCmd.PersistentFlags().StringP("github-team", "n", "", "Github team name")
	RootCmd.PersistentFlags().IntP("github-team-id", "i", 0, "Github team id")

	viper.BindPFlag("github_api_token", RootCmd.PersistentFlags().Lookup("github-api-token"))
	viper.BindPFlag("github_organization", RootCmd.PersistentFlags().Lookup("github-organization"))
	viper.BindPFlag("github_team", RootCmd.PersistentFlags().Lookup("github-team"))
	viper.BindPFlag("github_team_id", RootCmd.PersistentFlags().Lookup("github-team-id"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".github-authorized-keys") // name of config file (without extension)
	viper.AddConfigPath("$HOME")                   // adding home directory as first search path
	viper.AutomaticEnv()                           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
