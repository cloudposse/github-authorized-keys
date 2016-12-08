package cmd

import (
	"github.com/cloudposse/github-authorized-keys/key_storages"
)

func authorize(cfg config, userName string) (string, error) {
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
