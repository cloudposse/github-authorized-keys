package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("GithubKeys as backend storage", func() {
	var (
		validToken string
		validOrg string
		validTeamName string
		validTeamID int
		validUser string
	)

	BeforeEach(func() {
		validToken = viper.GetString("github_api_token")
		validOrg = viper.GetString("github_organization")
		validTeamName = viper.GetString("github_team")
		validTeamID = viper.GetInt("github_team_id")
		validUser = viper.GetString("github_user")
	})

	Describe("when github.com up", func() {
		var c *githubKeys

		BeforeEach(func() {
			c = newGithubKeys(validToken, validOrg, validTeamName, validTeamID)
		})

		Context("backend have valid value", func() {

			It("should return valid value", func() {

				keys, err := c.Get(validUser)

				Expect(err).To(BeNil())
				Expect(keys).NotTo(Equal(""))
			})
		})

		Context("backend does not have valid value", func() {
			It("should return empty value and valid error", func() {

				keys, err := c.Get("djahsjdhadafdsgfhdgahfjd")

				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(ErrStorageKeyNotFound))
				Expect(keys).To(Equal(""))
			})
		})
	})

	Describe("when github.com down", func() {

		var c *githubKeys

		BeforeEach(func() {
			c = newGithubKeys(validToken, validOrg, validTeamName, validTeamID)
			c.client.client.BaseURL.Host = "github.ru"
		})

		It("should return valid error and empty string", func() {
			keys, err := c.Get(validUser)

			Expect(err).NotTo(BeNil())
			Expect(keys).To(Equal(""))
		})
	})

})