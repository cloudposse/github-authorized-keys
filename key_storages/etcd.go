package keyStorages

import (
	"github.com/coreos/etcd/client"
	"time"
	"golang.org/x/net/context"
)

// ETCDCache - ETCD based key storage used as cache
type ETCDCache struct {
	client client.Client
	options *client.SetOptions
}

// Get - fetch value from key storage
func (c *ETCDCache) Get(key string) (value string, err error) {
	kapi := client.NewKeysAPI(c.client)
	resp, err := kapi.Get(context.Background(), key, nil)


	defer func() {
		if r := recover(); r != nil {
			value = ""
			err = ErrStorageConnectionFailed
		}
	}()

	switch err {
		case nil:
			value = resp.Node.Value
		default:
			value = ""
			if err.(client.Error).Code == 100 {
				err = ErrStorageKeyNotFound
			}

	}

	return
}

// Set - save value into key storage
func (c *ETCDCache) Set(key, value string) (err error) {
	kapi := client.NewKeysAPI(c.client)
	_, err = kapi.Set(context.Background(), key, value, c.options)


	if _, ok := err.(*client.ClusterError); ok {
		err = ErrStorageConnectionFailed
	}

	return
}

// Remove - remove value by key from key storage
func (c *ETCDCache) Remove(key string) (err error) {
	kapi := client.NewKeysAPI(c.client)
	_, err = kapi.Delete(context.Background(), key, nil)

	if _, ok := err.(*client.ClusterError); ok {
		err = ErrStorageConnectionFailed
	}

	return
}

// NewEtcdCache - constructor for etcd based key storage
func NewEtcdCache(endpoints []string, ttl time.Duration) (*ETCDCache, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	options := &client.SetOptions{TTL: ttl}

	return &ETCDCache{client: c, options: options}, err
}