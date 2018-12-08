// +build !windows,!darwin

package config

import (
	"os"
	"path/filepath"
)

const configHome = ".uping.toml"
const configSys = "config.toml"

const etc = "/etc/uping/"

func getConfigFileName() string {
	cfgPath := filepath.Join(os.Getenv("HOME"), configHome)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	cfgPath = filepath.Join(etc, configSys)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	return ""
}
