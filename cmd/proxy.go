package cmd

import (
	log "github.com/Sirupsen/logrus"
	"errors"
)

var (
	// ErrStorageKeyNotFound - returned when value is not fould in storage (source or fallback cache)
	ErrStorageKeyNotFound = errors.New("Storage: Key not found")

	// ErrStorageConnectionFailed - returned when there was connection error to storage (source or fallback cache)
	ErrStorageConnectionFailed =    errors.New("Storage: Connection failed")
)

type fallbackCache interface {
	Get(key string) (string, error)
	Set(key, value string) (error)
	Remove(key string) (error)
}

type source interface {
	Get(key string) (string, error)
}


type proxy struct {
	fallbackCache fallbackCache
	source        source
}

func (c *proxy) Get(name string) (value string, err error) {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "Get"})

	logger.Debugf("Backend lookup %v", name)

	value, err = c.lookupIn(c.source, name)

	switch err {
		case nil:
			logger.Debugf("Backend found %v", name)
			c.saveTo(c.fallbackCache, name, value)
			return

		case ErrStorageKeyNotFound:
			c.removeFrom(c.fallbackCache, name)
			return

		default:
			logger.Debug("Backend failed")
			logger.Debug("Fallback to cache")
			value, err = c.fallbackCache.Get(name)
			return
	}
}

func (c *proxy) lookupIn(storage source, name string) (string, error) {
	return storage.Get(name)
}

func (c *proxy) saveTo(storage fallbackCache, name, value string) error {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "saveTo"})
	logger.Debugf("Saving to cache %v: %v", name, value)
	return storage.Set(name, value)
}

func (c *proxy) removeFrom(storage fallbackCache, name string) error {
	logger := log.WithFields(log.Fields{"class": "Proxy", "method": "removeFrom"})
	logger.Debugf("Remove %v from cache", name)
	return storage.Remove(name)
}