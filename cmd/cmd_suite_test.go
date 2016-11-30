package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
	"github.com/spf13/viper"
)



func TestSuite(t *testing.T) {
	viper.SetConfigName(".github-authorized-keys-tests") // name of config file (without extension)
	viper.SetConfigType("yaml") // Set config format to yaml
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AddConfigPath("../")               // optionally look for config in the working directory
	viper.SetEnvPrefix("TEST")
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Github Authorized Keys Suite")
}
