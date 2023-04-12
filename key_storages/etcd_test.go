package keyStorages

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("ETCD", func() {

	const (
		validKey       = "TestKey"
		validValue     = "TestValue"
		testETCDPrefix = "/github-authorized-keys/tests"
	)

	var (
		endpoints []string
		ttl       time.Duration
	)

	BeforeEach(func() {
		ttl = 10 * time.Millisecond
	})

	Describe("with valid connection url", func() {
		BeforeEach(func() {
			endpoints = []string{}
			if etcd := viper.GetString("etcd_endpoint"); etcd != "" {
				endpoints = strings.Split(etcd, ",")
			}
		})

		Describe("constructor newEtcdCache()", func() {
			Context("call with valid connection url", func() {
				It("should return nil error", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCD_ENDPOINT to run this test")
					}
					_, err := NewEtcdCache(endpoints, testETCDPrefix, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCD_ENDPOINT to run this test")
					}

					client, _ := NewEtcdCache(endpoints, testETCDPrefix, ttl)

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
				client, _ = NewEtcdCache(endpoints, testETCDPrefix, ttl)

			})

			Context("call with existed key", func() {
				It("should return valid value and nil error ", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCD_ENDPOINT to run this test")
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
						Skip("Specify TEST_ETCD_ENDPOINT to run this test")
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
				client, _ = NewEtcdCache(endpoints, testETCDPrefix, ttl)
			})

			Context("call with removed existed key", func() {
				It("should return empty value and valid error ", func() {
					if len(endpoints) == 0 {
						Skip("Specify TEST_ETCD_ENDPOINT to run this test")
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
					_, err := NewEtcdCache(endpoints, testETCDPrefix, ttl)
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Set()", func() {
			Context(" key => value", func() {
				It("should return nil error", func() {
					client, _ := NewEtcdCache(endpoints, testETCDPrefix, ttl)

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
				client, _ = NewEtcdCache(endpoints, testETCDPrefix, ttl)
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
				client, _ = NewEtcdCache(endpoints, testETCDPrefix, ttl)
			})

			It("should return empty value and valid error ", func() {
				err := client.Remove(validKey)
				Expect(err).To(Equal(ErrStorageConnectionFailed))
			})
		})
	})
})
