package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/cloudposse/github-authorized-keys/api"
	"time"
)

var cfgFile string


type config struct {
	githubAPIToken string
	githubOrganization string
	githubTeamName string
	githubTeamID int

	etcdEndpoints 	[]string
	etcdTTL       string

	userGID    string
	userGroups []string
	userShell  string
	root   string
}

type flag struct {
	short string
	flagType string
	option string
	defaultValue interface{}
	description string
}

func (f *flag) flag() string {
	return strings.Replace(f.option, "_", "-", -1)
}


var flags = []flag{
flag{"t", "string",  "github_api_token",    "", 		       "Github API token    ( environment variable GITHUB_API_TOKEN could be used instead ) (read more https://github.com/blog/1509-personal-api-tokens)"},
flag{"o", "string",  "github_organization", "", 		       "Github organization ( environment variable GITHUB_ORGANIZATION could be used instead )"},
flag{"n", "string",  "github_team",         "", 		       "Github team name    ( environment variable GITHUB_TEAM could be used instead )"},
flag{"i", "int",     "github_team_id",       0, 		       "Github team id 	    ( environment variable GITHUB_TEAM_ID could be used instead )"},

flag{"g", "string",  "sync_users_gid",      "", 		       "Primary group id    ( environment variable SYNC_USERS_GID could be used instead )"},
flag{"G", "strings", "sync_users_groups",   []string{}, 	       "CSV groups name     ( environment variable SYNC_USERS_GROUPS could be used instead )"},
flag{"s", "string",  "sync_users_shell",    "/bin/bash",	       "User shell 	    ( environment variable SYNC_USERS_SHELL could be used instead )"},
flag{"r", "string",  "sync_users_root",     "/",		       "Root directory 	    ( environment variable SYNC_USERS_ROOT could be used instead )"},

flag{"e", "strings", "etcdctl_endpoint",    []string{},		       "CSV etcd endpoints  ( environment variable ETCDCTL_ENDPOINT could be used instead )"},
flag{"p", "string",  "etcdctl_prefix",      "/github-authorized-keys", "Path for etcd data  ( environment variable ETCDCTL_PREFIX could be used instead )"},
flag{"l", "int64",   "etcdctl_ttl",    	    ETCDTTLDefault,	       "ETCD value's ttl    ( environment variable ETCDCTL_TTL could be used instead )"},
}

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
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()

		cfg := config {
			githubAPIToken: 	viper.GetString("github_api_token"),
			githubOrganization: 	viper.GetString("github_organization"),
			githubTeamName:		viper.GetString("github_team"),
			githubTeamID:		viper.GetInt("github_team_id"),

			etcdEndpoints:	fixStringSlice(viper.GetString("etcdctl_endpoint")),
			etcdTTL: 	viper.GetString("etcdctl_ttl"),

			userGID: viper.GetString("sync_users_gid"),
			userGroups: fixStringSlice(viper.GetString("sync_users_groups")),
			userShell: viper.GetString("sync_users_shell"),
			root: viper.GetString("sync_users_root"),
		}


		

		router.GET("/authorize/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.JSON(200, gin.H{
				"message": "Authorize "+name,
			})
		})
		router.Run()
	},
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

	// Config file
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"Config file (default is $HOME/.github-authorized-keys.yaml)")

	for _, f := range flags {
		switch f.flagType {
		case "strings":
			RootCmd.Flags().StringSliceP(f.flag(), f.short, f.defaultValue.([]string), f.description)
			break
		case "int":
			RootCmd.Flags().IntP(f.flag(), f.short, f.defaultValue.(int), f.description)
			break
		case "int64":
			RootCmd.Flags().Int64P(f.flag(), f.short, f.defaultValue.(int64), f.description)
			break
		default:
			RootCmd.Flags().StringP(f.flag(), f.short, f.defaultValue.(string), f.description)
			break


		}
		viper.BindPFlag(f.option, RootCmd.Flags().Lookup(f.flag()))
	}
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


func fixStringSlice(s string) []string {
	result := []string{}
	if s != "" {
		result = strings.Split(s, ",")
	}
	return result
}



