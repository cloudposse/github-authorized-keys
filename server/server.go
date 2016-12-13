package server

import (
	"github.com/cloudposse/github-authorized-keys/config"
	"github.com/cloudposse/github-authorized-keys/key_storages"
	"github.com/gin-gonic/gin"
)

// Run - start http server
func Run(cfg config.Config) {

	router := gin.Default()

	router.GET("/user/:name/authorized_keys", func(c *gin.Context) {
		name := c.Param("name")
		key, err := authorize(cfg, name)
		if err == nil {
			c.String(200, "%v", key)
		} else {
			c.String(404, "")
		}
	})

	router.Run(cfg.Listen)
}

func authorize(cfg config.Config, userName string) (string, error) {
	var keys *keyStorages.Proxy

	sourceStorage := keyStorages.NewGithubKeys(cfg.GithubAPIToken, cfg.GithubOrganization, cfg.GithubTeamName, cfg.GithubTeamID)

	if len(cfg.EtcdEndpoints) > 0 {
		fallbackStorage, _ := keyStorages.NewEtcdCache(cfg.EtcdEndpoints, cfg.EtcdPrefix, cfg.EtcdTTL)
		keys = keyStorages.NewProxy(sourceStorage, fallbackStorage)
	} else {
		fallbackStorage := &keyStorages.NilStorage{}
		keys = keyStorages.NewProxy(sourceStorage, fallbackStorage)
	}
	return keys.Get(userName)
}
