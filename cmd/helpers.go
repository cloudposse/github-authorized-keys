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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

type flag struct {
	short        string
	flagType     string
	option       string
	defaultValue interface{}
	description  string
}

func (f *flag) flag() string {
	return strings.Replace(f.option, "_", "-", -1)
}

func createCmdFlags(cmd *cobra.Command, f flag) {
	switch f.flagType {
	case "strings":
		cmd.Flags().StringSliceP(f.flag(), f.short, f.defaultValue.([]string), f.description)
		break
	case "int":
		cmd.Flags().IntP(f.flag(), f.short, f.defaultValue.(int), f.description)
		break
	case "int64":
		cmd.Flags().Int64P(f.flag(), f.short, f.defaultValue.(int64), f.description)
		break
	case "bool":
		cmd.Flags().BoolP(f.flag(), f.short, f.defaultValue.(bool), f.description)
		break
	default:
		cmd.Flags().StringP(f.flag(), f.short, f.defaultValue.(string), f.description)
		break

	}
	viper.BindPFlag(f.option, cmd.Flags().Lookup(f.flag()))
}

func fixStringSlice(s string) []string {
	result := []string{}
	if s != "" {
		result = strings.Split(s, ",")
	}
	return result
}
