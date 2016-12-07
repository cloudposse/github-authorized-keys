package keyStorages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"github.com/spf13/viper"
	"strings"
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
		endpoints []string
		ttl time.Duration
	)

	BeforeEach(func() {
		validToken = viper.GetString("github_api_token")
		validOrg = viper.GetString("github_organization")
		validUser = viper.GetString("github_user")

		ttl = 10 * time.Millisecond
	})



	Describe("with valid connection url", func() {
		BeforeEach(func() {
			endpoints = []string{}
			if etcd := viper.GetString("etcdctl_endpoint"); etcd != "" {
				endpoints = strings.Split(etcd, ",")
			}
		})

		Describe("constructor newEtcdCache()", func() {
			Context("call with valid connection url", func() {
				It("should return nil error", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCDCTL_ENDPOINT to run this test")
					}
					_, err := NewEtcdCache(endpoints, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCDCTL_ENDPOINT to run this test")
					}

					client, _ := NewEtcdCache(endpoints, ttl)

					err := client.Set(validKey, validValue)

					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Get()", func() {
			var (
				client *ETCDCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(endpoints, ttl)

			})

			Context("call with existed key", func() {
				It("should return valid value and nil error ", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCDCTL_ENDPOINT to run this test")
					}

					client.Set(validKey, validValue)
					value, err := client.Get(validKey)
					Expect(err).To(BeNil())
					Expect(value).To(Equal(validValue))
				})
			})

			Context("call with existed key after ttl expired", func() {
				It("should return empty value and error ", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCDCTL_ENDPOINT to run this test")
					}

					client.Set(validKey, validValue)
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
				client *ETCDCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(endpoints, ttl)
			})

			Context("call with removed existed key", func() {
				It("should return empty value and valid error ", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCDCTL_ENDPOINT to run this test")
					}

					client.Set(validKey, validValue)

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
			endpoints = []string{"http://dasdsfafdsfa:2379"}
		})

		Describe("constructor newEtcdCache()", func() {
			Context("call with valid connection url", func() {
				It("should return nil error", func() {
					_, err := NewEtcdCache(endpoints, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					client, _ := NewEtcdCache(endpoints, ttl)

					err := client.Set(validKey, validValue)

					Expect(err).To(Equal(ErrStorageConnectionFailed))
				})
			})
		})

		Describe("Get()", func() {
			var (
				client *ETCDCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(endpoints, ttl)
			})

			It("should return empty value and valid error ", func() {
					value, err := client.Get(validKey)
					Expect(err).To(Equal(ErrStorageConnectionFailed))
					Expect(value).To(Equal(""))
			})
		})

		Describe("Remove()", func() {
			var (
				client *ETCDCache
			)

			BeforeEach(func() {
				client, _ = NewEtcdCache(endpoints, ttl)
			})

			It("should return empty value and valid error ", func() {
					err := client.Remove(validKey)
					Expect(err).To(Equal(ErrStorageConnectionFailed))
			})
		})
	})
})