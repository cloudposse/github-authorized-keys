package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("GithubClient", func() {
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

	Describe("getTeam()", func() {
		Context("call with valid token, org, team name and team id ", func() {
			It("should return nil error and valid team", func() {
				c := NewGithubClient(validToken, validOrg)
				team, err := c.GetTeam(validTeamName, validTeamID)

				Expect(err).To(BeNil())

				Expect(team).NotTo(BeNil())
				Expect(team.ID).NotTo(BeZero())
				Expect(*team.Name).NotTo(BeEmpty())
			})
		})

		Context("call with invalid team name AND valid token, org, team id", func() {
			It("should return nil error and valid team", func() {
				c := NewGithubClient(validToken, validOrg)
				team, err := c.GetTeam("dasdasd", validTeamID)

				Expect(err).To(BeNil())

				Expect(team).NotTo(BeNil())
				Expect(team.ID).NotTo(BeZero())
				Expect(*team.Name).NotTo(BeEmpty())
			})
		})

		Context("call with invalid team name && team id AND valid token, org", func() {
			It("should return valid error and nil team", func() {
				c := NewGithubClient(validToken, validOrg)
				team, err := c.GetTeam("dasdasd", 0)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Team with such name or id not found"))

				Expect(team).To(BeNil())
			})
		})

		Context("call with invalid token AND valid org, team name, team id", func() {
			It("should return valid error and nil team", func() {
				c := NewGithubClient("11111111111111111111111111", validOrg)
				team, err := c.GetTeam(validTeamName, validTeamID)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Access denied"))

				Expect(team).To(BeNil())
			})
		})

		Context("call with invalid org AND valid token, team name, team id", func() {
			It("should return valid error and nil team", func() {
				c := NewGithubClient(validToken, "dsadsad")
				team, err := c.GetTeam(validTeamName, validTeamID)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Access denied"))

				Expect(team).To(BeNil())
			})
		})

	})

	Describe("isTeamMember()", func() {
		Context("call with user that is member of the team", func() {
			It("should return nil error and true value", func() {
				c := NewGithubClient(validToken, validOrg)
				team, _ := c.GetTeam(validTeamName, validTeamID)
				isMember, err := c.IsTeamMember(validUser, team)

				Expect(err).To(BeNil())

				Expect(isMember).To(BeTrue())
			})
		})

		Context("call with user that is not member of the team", func() {
			It("should return nil error and false value", func() {
				c := NewGithubClient(validToken, validOrg)
				team, _ := c.GetTeam(validTeamName, validTeamID)
				isMember, err := c.IsTeamMember("dasda", team)

				Expect(err).To(BeNil())

				Expect(isMember).To(BeFalse())
			})
		})

	})

	Describe("getUser()", func() {
		Context("call with valid user", func() {
			It("should return nil error and not nil user", func() {
				c := NewGithubClient(validToken, validOrg)
				user, err := c.getUser(validUser)

				Expect(err).To(BeNil())

				Expect(user).NotTo(BeNil())
			})
		})

		Context("call with invalid user", func() {
			It("should return error and nil user", func() {
				c := NewGithubClient(validToken, validOrg)
				user, err := c.getUser("dasdddds232dasdas")

				Expect(err).NotTo(BeNil())

				Expect(user).To(BeNil())
			})
		})

	})

	Describe("getKeys()", func() {
		Context("call with valid user", func() {
			It("should return nil error and no empty list of keys", func() {
				c := NewGithubClient(validToken, validOrg)
				user, _ := c.getUser(validUser)
				keys, _, err := c.GetKeys(user.Login)

				Expect(err).To(BeNil())

				Expect(len(keys) > 0).To(BeTrue())
			})
		})
	})

	Describe("getTeamMembers()", func() {
		Context("call with valid team", func() {
			It("should return nil error and no empty list of members", func() {
				c := NewGithubClient(validToken, validOrg)
				team, _ := c.GetTeam(validTeamName, validTeamID)

				members, err := c.GetTeamMembers(team)

				Expect(err).To(BeNil())

				Expect(len(members) > 0).To(BeTrue())
			})
		})
	})
})