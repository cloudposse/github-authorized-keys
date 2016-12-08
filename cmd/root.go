package cmd

import (
	"fmt"
	"os"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var cfgFile string

type config struct {
	GithubAPIToken     string
	GithubOrganization string
	GithubTeamName     string
	GithubTeamID       int

	EtcdEndpoints []string
	EtcdTTL       time.Duration
	EtcdPrefix    string

	UserGID    string
	UserGroups []string
	UserShell  string
	Root       string
}

type flag struct {
	short        string
	flagType     string
	option       string
	defaultValue interface{}
	description  string
}

func (f *flag) flag() string {
	return strings.Replace(f.option, "_", "-", -1)
}

// ETCDTTLDefault - default ttl - 1day in seconds = 24 hours * 60 minutes * 60 seconds
const ETCDTTLDefault = int64(24 * 60 * 60)

var flags = []flag{
	{"t", "string", "github_api_token", "", "Github API token    ( environment variable GITHUB_API_TOKEN could be used instead ) (read more https://github.com/blog/1509-personal-api-tokens)"},
	{"o", "string", "github_organization", "", "Github organization ( environment variable GITHUB_ORGANIZATION could be used instead )"},
	{"n", "string", "github_team", "", "Github team name    ( environment variable GITHUB_TEAM could be used instead )"},
	{"i", "int", "github_team_id", 0, "Github team id 	    ( environment variable GITHUB_TEAM_ID could be used instead )"},

	{"g", "string", "sync_users_gid", "", "Primary group id    ( environment variable SYNC_USERS_GID could be used instead )"},
	{"G", "strings", "sync_users_groups", []string{}, "CSV groups name     ( environment variable SYNC_USERS_GROUPS could be used instead )"},
	{"s", "string", "sync_users_shell", "/bin/bash", "User shell 	    ( environment variable SYNC_USERS_SHELL could be used instead )"},
	{"r", "string", "sync_users_root", "/", "Root directory 	    ( environment variable SYNC_USERS_ROOT could be used instead )"},

	{"e", "strings", "etcdctl_endpoint", []string{}, "CSV etcd endpoints  ( environment variable ETCDCTL_ENDPOINT could be used instead )"},
	{"p", "string", "etcdctl_prefix", "/github-authorized-keys", "Path for etcd data  ( environment variable ETCDCTL_PREFIX could be used instead )"},
	{"l", "int64", "etcdctl_ttl", ETCDTTLDefault, "ETCD value's ttl    ( environment variable ETCDCTL_TTL could be used instead )"},
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
	RunE: func(cmd *cobra.Command, args []string) error {
		etcdTTL, err := time.ParseDuration(viper.GetString("etcdctl_ttl") + "s")

		if err != nil {
			return err
		}

		cfg := config{
			GithubAPIToken:     viper.GetString("github_api_token"),
			GithubOrganization: viper.GetString("github_organization"),
			GithubTeamName:     viper.GetString("github_team"),
			GithubTeamID:       viper.GetInt("github_team_id"),

			EtcdEndpoints: fixStringSlice(viper.GetString("etcdctl_endpoint")),
			EtcdTTL:       etcdTTL,

			UserGID:    viper.GetString("sync_users_gid"),
			UserGroups: fixStringSlice(viper.GetString("sync_users_groups")),
			UserShell:  viper.GetString("sync_users_shell"),
			Root:       viper.GetString("sync_users_root"),
		}

		err = validation.StructRules{}.
			Add("GithubAPIToken", validation.Required.Error("is required")).
			Add("GithubOrganization", validation.Required.Error("is required")).
			Add("EtcdEndpoints", is.URL).
			/*		// Should be valid duration in seconds
					Add("etcdTTL", func(value string) error {
							_, err := time.ParseDuration(value + "s")
							return err
					}).*/
			// performs validation
			Validate(cfg)

		if err != nil {
			return err
		}

		// Validate Github Team exists
		if cfg.GithubTeamName == "" && cfg.GithubTeamID == 0 {
			return errors.New("Team name or Team id should be specified")
		}

		router := gin.Default()

		router.GET("/authorize/:name", func(c *gin.Context) {
			name := c.Param("name")
			key, err := authorize(cfg, name)
			if err == nil {
				c.String(200, "%v", key)
			} else {
				c.String(404, "")
			}
		})

		router.Run()


		
		return nil
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
