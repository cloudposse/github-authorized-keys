package keyStorages

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

var (
	// ErrStorageKeyNotFound - returned when value is not found in storage (source or fallback cache)
	ErrStorageKeyNotFound = errors.New("storage: Key not found")

	// ErrStorageConnectionFailed - returned when there was connection error to storage (source or fallback cache)
	ErrStorageConnectionFailed = errors.New("storage: Connection failed")
)

type fallbackCache interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Remove(key string) error
}

type source interface {
	Get(key string) (string, error)
}

// Proxy - key storage fallback proxy.
// Always deal with source key storage first, and sync values with fallback cache storage
// If source key storage is unavailable fallback to cache storage
type Proxy struct {
	fallbackCache fallbackCache
	source        source
}

// Get - fetch value from key storage
func (c *Proxy) Get(name string) (value string, err error) {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "Get"})

	logger.Debugf("Backend lookup %v", name)

	value, err = c.lookupIn(c.source, name)

	switch err {
	case nil:
		logger.Debugf("Backend found %v", name)
		err = c.saveTo(c.fallbackCache, name, value)
		if err != nil {
			return "", err
		}
		return

	case ErrStorageKeyNotFound:
		err = c.removeFrom(c.fallbackCache, name)
		if err != nil {
			return "", err
		}
		return

	default:
		logger.Debug("Backend failed")
		logger.Debug("Fallback to cache")
		value, err = c.fallbackCache.Get(name)
		return
	}
}

func (c *Proxy) lookupIn(storage source, name string) (string, error) {
	return storage.Get(name)
}

func (c *Proxy) saveTo(storage fallbackCache, name, value string) error {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "saveTo"})
	logger.Debugf("Saving to cache %v: %v", name, value)
	return storage.Set(name, value)
}

func (c *Proxy) removeFrom(storage fallbackCache, name string) error {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "removeFrom"})
	logger.Debugf("Remove %v from cache", name)
	return storage.Remove(name)
}

// NewProxy - constructor to create Proxy object
func NewProxy(source source, fallbackCache fallbackCache) *Proxy {
	return &Proxy{source: source, fallbackCache: fallbackCache}
}
