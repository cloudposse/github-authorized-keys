package keyStorages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CacheMap struct {
	storage *map[string]string
}

func (c *CacheMap) Get(name string) (string, error) {
	val, ok := (*c.storage)[name]
	if ! ok {
		return "", ErrStorageKeyNotFound
	}
	return	val, nil
}

func (c *CacheMap) Set(name, value string) (error) {
	(*c.storage)[name] = value
	return nil
}

func (c *CacheMap) Remove(name string) (error) {
	delete(*c.storage, name)
	return nil
}

type BackendMap struct {
	storage *map[string]string
}

func (c *BackendMap) Get(name string) (string, error) {
	val, ok := (*c.storage)[name]
	if ! ok {
		return "", ErrStorageKeyNotFound
	}
	return	val, nil
}

type BackendFail struct {}

func (c *BackendFail) Get(name string) (string, error) {
	return "", ErrStorageConnectionFailed
}

var _ = Describe("Proxy", func() {
	var (
		cacheStorage map[string]string
		backendStorage map[string]string
		proxyStorage Proxy
	)

	Context("backend have valid value", func() {
		BeforeEach(func() {
			cacheStorage = map[string]string{}
			backendStorage = map[string]string{}

			proxyStorage = Proxy{
				fallbackCache:   &CacheMap{storage: &cacheStorage},
				source: &BackendMap{storage: &backendStorage},
			}

			backendStorage["goruha"] = "TestValue"
		})

		It("should return valid value", func() {
			value, err := proxyStorage.Get("goruha")

			Expect(err).To(BeNil())
			Expect(value).To(Equal("TestValue"))
		})

		It("should save value to cache", func() {
			proxyStorage.Get("goruha")

			value, ok := cacheStorage["goruha"]

			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("TestValue"))
		})
	})

	Context("backend does not valid value, but cache do", func() {
		BeforeEach(func() {
			cacheStorage = map[string]string{}
			backendStorage = map[string]string{}

			proxyStorage = Proxy{
				fallbackCache:   &CacheMap{storage: &cacheStorage},
				source: &BackendMap{storage: &backendStorage},
			}

			cacheStorage["goruha"] = "TestValue"
		})

		It("should return empty value and error", func() {
			value, err := proxyStorage.Get("goruha")

			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(ErrStorageKeyNotFound))
			Expect(value).To(Equal(""))
		})

		It("should remove value from cache", func() {
			proxyStorage.Get("goruha")

			value, ok := cacheStorage["goruha"]

			Expect(ok).To(BeFalse())
			Expect(value).To(Equal(""))
		})
	})

	Context("neither backend nor cache store valid value", func() {
		BeforeEach(func() {
			cacheStorage = map[string]string{}
			backendStorage = map[string]string{}

			proxyStorage = Proxy{
				fallbackCache:   &CacheMap{storage: &cacheStorage},
				source: &BackendMap{storage: &backendStorage},
			}
		})

		It("should return empty value and error", func() {
			value, err := proxyStorage.Get("goruha")

			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(ErrStorageKeyNotFound))
			Expect(value).To(Equal(""))
		})
	})

	Context("backend failed to be connected, but cache do have correct value", func() {
		BeforeEach(func() {
			cacheStorage = map[string]string{}

			proxyStorage = Proxy{
				fallbackCache:   &CacheMap{storage: &cacheStorage},
				source: &BackendFail{},
			}

			cacheStorage["goruha"] = "TestValue"
		})

		It("should return valid value and nil error", func() {
			value, err := proxyStorage.Get("goruha")

			Expect(err).To(BeNil())
			Expect(value).To(Equal("TestValue"))
		})
	})
})