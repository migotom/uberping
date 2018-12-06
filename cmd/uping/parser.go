package main

import (
	"log"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/migotom/uberping/internal/driver"
	"github.com/migotom/uberping/internal/schema"
	"github.com/migotom/uberping/internal/worker"
)

func configParser(arguments map[string]interface{}, appConfig *schema.GeneralConfig) (schema.HostsLoader, worker.ResultsSaver) {
	var hostsLoader schema.HostsLoader
	var resultsSaver worker.ResultsSaver

	// Parse arguments

	// Load config
	if apiConfigFile, ok := arguments["-C"].(string); ok {
		if _, err := toml.DecodeFile(apiConfigFile, &appConfig); err != nil {
			log.Fatal(err)
		}
	}

	// Override config by args
	appConfig.Verbose = !arguments["-s"].(bool)
	appConfig.Grouped = arguments["-g"].(bool)

	appConfig.Ping.Privileged = true
	proto := arguments["-p"].(string)
	if proto == "udp" {
		appConfig.Ping.Privileged = false
	}

	appConfig.Ping.Interval = time.Duration(1) * time.Second
	if interval, ok := arguments["-i"].(string); ok {
		if interval, err := time.ParseDuration(interval); err == nil {
			appConfig.Ping.Interval = interval
		}
	}

	appConfig.Ping.Count = 4
	if count, ok := arguments["-c"].(string); ok {
		if count, err := strconv.ParseInt(count, 10, 64); err == nil {
			appConfig.Ping.Count = int(count)
		}
	}

	appConfig.Ping.Timeout = time.Duration(int(appConfig.Ping.Count)) * time.Second
	if timeout, ok := arguments["-t"].(string); ok {
		if timeout, err := time.ParseDuration(timeout); err == nil {
			appConfig.Ping.Timeout = timeout
		}
	}

	appConfig.Workers = 4
	if workers, ok := arguments["-w"].(string); ok {
		if workers, err := strconv.ParseInt(workers, 10, 64); err == nil {
			appConfig.Workers = int(workers)
		}
	}

	if hosts, ok := arguments["<hosts>"].([]string); ok {
		hostsLoader = func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.ArgvLoadHosts(parser, hosts)
		}
	}

	if file, ok := arguments["<file-in>"].(string); ok {
		hostsLoader = func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.FileLoadHosts(parser, file)
		}
	}

	if api := arguments["--source-api"].(bool); api {
		hostsLoader = func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.APILoadHosts(parser, &appConfig.API)
		}
	}

	if api := arguments["--out-api"].(bool); api {
		resultsSaver = func(pingResult schema.PingResult) error {
			return driver.APISavePingResult(pingResult, &appConfig.API)
		}
	}

	if file, ok := arguments["<file-out>"].(string); ok {
		resultsSaver = func(pingResult schema.PingResult) error {
			return driver.FileSavePingResult(pingResult, file)
		}
	}

	if resultsSaver == nil {
		resultsSaver = func(pingResult schema.PingResult) error {
			return driver.StdoutPingResult(pingResult)
		}
	}

	return hostsLoader, resultsSaver
}
