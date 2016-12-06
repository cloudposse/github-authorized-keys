package key_storages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"github.com/spf13/viper"
)

var _ = Describe("ETCD", func() {
	const (
		validKey = "TestKey"
		validValue = "TestValue"
	)

	var (
		validToken string
		validOrg string
		validUser string
	)

	var (
		gateways []string
		ttl time.Duration
	)

	BeforeEach(func() {
		viper.SetDefault("etcd", "http://etcd:2379")

		validToken = viper.GetString("github_api_token")
		validOrg = viper.GetString("github_organization")
		validUser = viper.GetString("github_user")

		ttl = 10 * time.Millisecond
	})

	Describe("with valid connection url", func() {
		BeforeEach(func() {
			gateways = []string{viper.GetString("etcd")}
		})

		Describe("constructor newEtcdCache()", func() {
			Context("call with valid connection url", func() {
				It("should return nil error", func() {
					_, err := NewEtcdCache(gateways, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					client, _ := NewEtcdCache(gateways, ttl)

					err := client.Set(validKey, validValue)

					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Get()", func() {
			var (
				client *etcdCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(gateways, ttl)
				client.Set(validKey, validValue)
			})

			Context("call with existed key", func() {
				It("should return valid value and nil error ", func() {
					value, err := client.Get(validKey)
					Expect(err).To(BeNil())
					Expect(value).To(Equal(validValue))
				})
			})

			Context("call with existed key after ttl expired", func() {
				It("should return empty value and error ", func() {
					time.Sleep(time.Second)
					value, err := client.Get(validKey)

					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(ErrStorageKeyNotFound))
					Expect(value).To(Equal(""))
				})
			})
		})

		Describe("Remove()", func() {
			var (
				client *etcdCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(gateways, ttl)
				client.Set(validKey, validValue)
			})

			Context("call with removed existed key", func() {
				It("should return empty value and valid error ", func() {
					client.Remove(validKey)
					value, err := client.Get(validKey)
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(ErrStorageKeyNotFound))
					Expect(value).To(Equal(""))
				})
			})
		})
	})

	Describe("with invalid connection url", func() {
		BeforeEach(func() {
			gateways = []string{"http://dasdsfafdsfa:2379"}
		})

		Describe("constructor newEtcdCache()", func() {
			Context("call with valid connection url", func() {
				It("should return nil error", func() {
					_, err := NewEtcdCache(gateways, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					client, _ := NewEtcdCache(gateways, ttl)

					err := client.Set(validKey, validValue)

					Expect(err).To(Equal(ErrStorageConnectionFailed))
				})
			})
		})

		Describe("Get()", func() {
			var (
				client *etcdCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(gateways, ttl)
			})

			It("should return empty value and valid error ", func() {
					value, err := client.Get(validKey)
					Expect(err).To(Equal(ErrStorageConnectionFailed))
					Expect(value).To(Equal(""))
			})
		})

		Describe("Remove()", func() {
			var (
				client *etcdCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(gateways, ttl)
			})

			It("should return empty value and valid error ", func() {
					err := client.Remove(validKey)
					Expect(err).To(Equal(ErrStorageConnectionFailed))
			})
		})
	})
})