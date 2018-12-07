package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	docopt "github.com/docopt/docopt-go"
	"github.com/migotom/uberping/internal/schema"
	"github.com/migotom/uberping/internal/worker"
)

var usage = `Uberping.

Usage:
  uping [options] [<hosts>...]
  uping -h | --help
  uping --version

Options:
  -C <config-file>         Use configuration file, eg. API endpoints, secrets, etc...
  -s                       Be silent and don't print output to stdout
  -g                       Print grouped results
  -p udp|icmp              Set type of ping packet, unprivileged udp or privileged icmp [default: icmp]
  -f				       Use fallback mode, uping will try to use next ping mode if selected by -p failed
  -c <count>               Number of pings to perform [default: 4]
  -i <ping-interval>       Interval between pings, eg. -i 1s, -i 100ms [default: 1s]
  -t <host-timeout>        Timeout before probing one host terminates, regardless of how many pings perfomed, eg. -t 1s, -t 100ms [default: <count> * 1s]
  -w <workers>             Number of paraller workers to run [default: 4]
  --source-api             Load hosts using external API configured by -C <config-file>
  --source-file <file-in>  Load hosts from file <file-in>
  --out-api                Save tests results using external API configured by -C <config-file>
  --out-file <file-out>    Save tests results to file <file-out>
`

const version = "0.1"

func main() {
	var Hosts schema.Hosts

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], version)
	fmt.Println(arguments)

	appConfig := schema.GeneralConfig{}
	hostsLoaders, resultsSavers := configParser(arguments, &appConfig)

	for _, hostsLoader := range hostsLoaders {
		if err := Hosts.Add(hostsLoader); err != nil {
			log.Fatal(err)
		}
	}

	// Create workers pool
	jobs := make(chan schema.Host, appConfig.Workers)
	appConfig.Results = make(chan schema.PingResult, appConfig.Workers*2)

	var wgPinger sync.WaitGroup
	var wgWriter sync.WaitGroup

	wgWriter.Add(1)
	go worker.Saver(0, appConfig, resultsSavers, &wgWriter)

	for i := 0; i < appConfig.Workers; i++ {
		wgPinger.Add(1)
		go worker.Pinger(i, appConfig, jobs, &wgPinger)
	}

	// Assign jobs (hosts to test)
	for _, host := range Hosts.Get() {
		jobs <- host
	}

	close(jobs)
	wgPinger.Wait()

	close(appConfig.Results)
	wgWriter.Wait()
}