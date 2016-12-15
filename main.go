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

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cloudposse/github-authorized-keys/cmd"
	"os"
	"github.com/spf13/viper"
)

func main() {
	LoggerInit()
	cmd.Execute()
}

// LoggerInit - Initialize logger configuration used for cli
func LoggerInit() {
	viper.SetDefault("LOG_LEVEL", "info")

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	loglevel :=  viper.GetString("LOG_LEVEL")
	switch loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break

	case "info":
	default:
		log.SetLevel(log.InfoLevel)
		break
	}
}
