/*
 * Github Authorized Keys - Use GitHub teams to manage system user accounts and authorized_keys
 *
 * Copyright 2016 Cloud Posse, LLC <hello@cloudposse.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudposse/github-authorized-keys/config"
	keyStorages "github.com/cloudposse/github-authorized-keys/key_storages"
	"github.com/gin-gonic/gin"
)

// Run - start http server
func Run(cfg config.Config) {

	router := gin.Default()

	router.GET("/user/:name/authorized_keys", func(c *gin.Context) {
		name := c.Params.ByName("name")
		name = strings.ToLower(name)
		key, err := authorize(cfg, name)
		if err == nil {
			c.String(200, "%v", key)
		} else {
			c.String(404, "")
		}
	})

	err := router.Run(cfg.Listen)
	if err != nil {
		log.Error(err)
	}
}

func authorize(cfg config.Config, userName string) (string, error) {
	var keys *keyStorages.Proxy

	sourceStorage := keyStorages.NewGithubKeys(cfg.GithubAPIToken, cfg.GithubOrganization, cfg.GithubURL, cfg.GithubTeamName, cfg.GithubTeamID)

	if len(cfg.EtcdEndpoints) > 0 {
		fallbackStorage, _ := keyStorages.NewEtcdCache(cfg.EtcdEndpoints, cfg.EtcdPrefix, cfg.EtcdTTL)
		keys = keyStorages.NewProxy(sourceStorage, fallbackStorage)
	} else {
		fallbackStorage := &keyStorages.NilStorage{}
		keys = keyStorages.NewProxy(sourceStorage, fallbackStorage)
	}
	return keys.Get(userName)
}
