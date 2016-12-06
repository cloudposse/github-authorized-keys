package key_storages

import (
	"github.com/coreos/etcd/client"
	"time"
	"golang.org/x/net/context"
)

type etcdCache struct {
	client client.Client
	options *client.SetOptions
}

func (c *etcdCache) Get(name string) (value string, err error) {
	kapi := client.NewKeysAPI(c.client)
	resp, err := kapi.Get(context.Background(), name, nil)


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

func (c *etcdCache) Set(name, value string) (err error) {
	kapi := client.NewKeysAPI(c.client)
	_, err = kapi.Set(context.Background(), name, value, c.options)


	if _, ok := err.(*client.ClusterError); ok {
		err = ErrStorageConnectionFailed
	}

	return
}

func (c *etcdCache) Remove(name string) (err error) {
	kapi := client.NewKeysAPI(c.client)
	_, err = kapi.Delete(context.Background(), name, nil)

	if _, ok := err.(*client.ClusterError); ok {
		err = ErrStorageConnectionFailed
	}

	return
}

func NewEtcdCache(gateways []string, ttl time.Duration) (*etcdCache, error) {
	cfg := client.Config{
		Endpoints:               gateways,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	options := &client.SetOptions{TTL: ttl}

	return &etcdCache{client: c, options: options}, err
}