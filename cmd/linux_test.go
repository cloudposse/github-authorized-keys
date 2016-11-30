package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"os/user"
	"strconv"
)

var _ = Describe("Linux api", func() {
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

	Describe("User exists", func() {
		Context("for no existed user", func() {
			It("should return false", func() {
				isExists := linuxUserExists("testdsadasfsa")
				Expect(isExists).To(BeFalse())
			})
		})

		Context("for existed user", func() {
			It("should return true", func() {
				isExists := linuxUserExists("root")
				Expect(isExists).To(BeTrue())
			})
		})
	})

	Describe("Create user", func() {
		Context("without GID", func() {
			var (
				userName linuxUser
			)

			BeforeEach(func() {
				userName = linuxUser{Gid: "", Name: "test", Shell: "/bin/bash", Groups: []string{"wheel", "root"}}
			})

			AfterEach(func() {
				linuxUserDelete(userName)
			})

			It("should create valid user", func() {
				err := linuxUserCreate(userName)

				Expect(err).To(BeNil())

				osUser, _ := user.Lookup(userName.Name)

				Expect(osUser.Username).To(Equal(userName.Name))

				value, _ := strconv.ParseInt(osUser.Gid, 10, 64);
				Expect(value > 0).To(BeTrue())

				gids, _:= osUser.GroupIds()

				for _, group := range userName.Groups {
					linuxGroup, err := user.LookupGroup(group)
					Expect(err).To(BeNil())
					Expect(gids).To(ContainElement(string(linuxGroup.Gid)))
				}

				shell := linuxUserShell(userName.Name)

				Expect(shell).To(Equal(userName.Shell))
			})
		})

		Context("with GID", func() {
			var (
				userName linuxUser
			)

			BeforeEach(func() {
				userName = linuxUser{Gid: "42", Name: "test", Shell: "/bin/bash", Groups: []string{"root"}}
			})

			AfterEach(func() {
				linuxUserDelete(userName)
			})

			It("should create valid user", func() {
				err := linuxUserCreate(userName)

				Expect(err).To(BeNil())

				osUser, _ := user.Lookup(userName.Name)

				Expect(osUser.Username).To(Equal(userName.Name))

				Expect(string(osUser.Gid)).To(Equal(userName.Gid))

				gids, _:= osUser.GroupIds()

				for _, group := range userName.Groups {
					linuxGroup, err := user.LookupGroup(group)
					Expect(err).To(BeNil())
					Expect(gids).To(ContainElement(string(linuxGroup.Gid)))
				}

				shell := linuxUserShell(userName.Name)

				Expect(shell).To(Equal(userName.Shell))
			})
		})
	})


	Describe("Group exists", func() {
		Context("for no existed group", func() {
			It("should return false", func() {
				isExists := linuxGroupExists("testdsadasfsa")
				Expect(isExists).To(BeFalse())
			})
		})

		Context("for existed group", func() {
			It("should return true", func() {
				isExists := linuxGroupExists("wheel")
				Expect(isExists).To(BeTrue())
			})
		})
	})


	Describe("Group exists by id", func() {
		Context("for no existed group", func() {
			It("should return false", func() {
				isExists := linuxGroupExistsByID("43")
				Expect(isExists).To(BeFalse())
			})
		})

		Context("for existed group", func() {
			It("should return true", func() {
				isExists := linuxGroupExistsByID("42")
				Expect(isExists).To(BeTrue())
			})
		})
	})

	Describe("Get User shell", func() {
		Context("for existed user", func() {
			It("should return /bin/ash", func() {
				shell := linuxUserShell("root")
				Expect(shell).To(Equal("/bin/ash"))
			})
		})
	})

})
