package main

import (
	"log"
	"strconv"
	"time"

	"github.com/migotom/uberping/internal/driver"
	"github.com/migotom/uberping/internal/schema"
	"github.com/migotom/uberping/internal/schema/config"
	"github.com/migotom/uberping/internal/worker"
)

func configParser(arguments map[string]interface{}, appConfig *schema.GeneralConfig) ([]schema.HostsLoader, []worker.ResultsSaver, []schema.HostsCleaner) {
	var hostsLoaders []schema.HostsLoader
	var resultsSavers []worker.ResultsSaver
	var cleaners []schema.HostsCleaner

	// Parse arguments

	// Load config
	var apiConfigFile string
	apiConfigFile, _ = arguments["-C"].(string)
	if err := config.LoadConfigFile(appConfig, apiConfigFile); err != nil {
		log.Fatal(err)
	}

	// Override config by args
	appConfig.Verbose = !arguments["-s"].(bool)
	appConfig.Grouped = arguments["-g"].(bool)

	if mode, ok := arguments["--mode"].(string); ok {
		appConfig.Probe.Mode = mode
	}
	switch appConfig.Probe.Mode {
	case "udp":
		appConfig.Probe.Worker = worker.Pinger
		appConfig.Probe.Privileged = false
	case "icmp":
		appConfig.Probe.Worker = worker.Pinger
		appConfig.Probe.Privileged = true
	case "netcat":
		appConfig.Probe.Worker = worker.Netcat
		appConfig.Probe.Mode = "netcat"
	case "":
		appConfig.Probe.Worker = worker.Pinger
		appConfig.Probe.Privileged = true
	default:
		log.Fatalln("Unsupported protocol.")
	}

	if defaultPort, ok := arguments["-P"].(string); ok {
		if defaultPort, err := strconv.ParseInt(defaultPort, 10, 64); err == nil {
			appConfig.Probe.DefaultPort = int(defaultPort)
		}
	}

	if appConfig.Verbose {
		resultsSavers = append(resultsSavers, func(pingResult schema.ProbeResult) error {
			return driver.StdoutPingResult(pingResult)
		})
	}

	if interval, ok := arguments["-d"].(string); ok {
		if interval, err := time.ParseDuration(interval); err == nil {
			appConfig.TestsInterval.Duration = interval
		}
	}

	if appConfig.Probe.Interval.Duration.Seconds() == 0 {
		appConfig.Probe.Interval.Duration = time.Duration(1) * time.Second
	}
	if interval, ok := arguments["-i"].(string); ok {
		if interval, err := time.ParseDuration(interval); err == nil {
			appConfig.Probe.Interval.Duration = interval
		}
	}

	if appConfig.Probe.Count == 0 {
		appConfig.Probe.Count = 4
	}
	if count, ok := arguments["-c"].(string); ok {
		if count, err := strconv.ParseInt(count, 10, 64); err == nil {
			appConfig.Probe.Count = int(count)
		}
	}

	if appConfig.Probe.Timeout.Duration.Seconds() == 0 {
		appConfig.Probe.Timeout.Duration = time.Duration(int(appConfig.Probe.Count)) * time.Second
	}
	if timeout, ok := arguments["-t"].(string); ok {
		if timeout, err := time.ParseDuration(timeout); err == nil {
			appConfig.Probe.Timeout.Duration = timeout
		}
	}

	if appConfig.Workers == 0 {
		appConfig.Workers = 4
	}
	if workers, ok := arguments["-w"].(string); ok {
		if workers, err := strconv.ParseInt(workers, 10, 64); err == nil {
			appConfig.Workers = int(workers)
		}
	}

	if hosts, ok := arguments["<hosts>"].([]string); ok {
		hostsLoaders = append(hostsLoaders, func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.ArgvLoadHosts(parser, hosts)
		})
	}

	if file, ok := arguments["--source-file"].(string); ok {
		hostsLoaders = append(hostsLoaders, func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.FileLoadHosts(parser, file)
		})
	}

	if db := arguments["--source-db"].(bool); db {
		hostsLoaders = append(hostsLoaders, func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.DBSqlLoadHosts(parser, &appConfig.DB)
		})
		cleaners = append(cleaners, func() {
			driver.DBCleaner(&appConfig.DB)
		})
	}

	if api := arguments["--source-api"].(bool); api {
		hostsLoaders = append(hostsLoaders, func(parser schema.HostParser) ([]schema.Host, error) {
			return driver.APILoadHosts(parser, &appConfig.API)
		})
	}

	if api := arguments["--out-api"].(bool); api {
		resultsSavers = append(resultsSavers, func(pingResult schema.ProbeResult) error {
			return driver.APISavePingResult(pingResult, &appConfig.API)
		})
	}

	if file, ok := arguments["--out-file"].(string); ok {
		resultsSavers = append(resultsSavers, func(pingResult schema.ProbeResult) error {
			return driver.FileSavePingResult(pingResult, file)
		})
	}

	if db := arguments["--out-db"].(bool); db {
		resultsSavers = append(resultsSavers, func(pingResult schema.ProbeResult) error {
			return driver.DBSqlSavePingResult(pingResult, &appConfig.DB)
		})
		cleaners = append(cleaners, func() {
			driver.DBCleaner(&appConfig.DB)
		})
	}

	return hostsLoaders, resultsSavers, cleaners
}
