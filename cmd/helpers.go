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
