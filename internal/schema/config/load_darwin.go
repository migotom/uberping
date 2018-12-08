// +build darwin

package config

import (
	"os"
	"path/filepath"
)

const configName = "config.toml"
const configHome = ".uping.toml"
const appSupport = "/Library/Application Support/Uping"

func getConfigFileName() string {
	cfgPath := filepath.Join(os.Getenv("HOME"), appSupport, configName)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	cfgPath = filepath.Join(os.Getenv("HOME"), configHome)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	const appSupport = "/Library/Application Support"
	cfgPath = filepath.Join(appSupport, configName)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	return ""
}
