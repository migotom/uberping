package main

import (
	"log"
	"os"
	"sync"
	"time"

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
  --mode udp|icmp|netcat   Set type of probe operation: ping with unprivileged udp, icmp or try to connect using tcp port (default: icmp)
  -d <tests-interval>      Interval between tests, if provided uping will perform tests indefinitely, e.g. every -I 1m, -I 1m30s, -I 1h30m10s
  -C <config-file>         Use configuration file, eg. API endpoints, secrets, etc...
  -s                       Be silent and don't print output to stdout, only errors to stderr
  -g                       Print grouped results
  -P <default-port>        In case of netcat mode use <default-port> for hosts without explicitly specified port, e.g. -p 8080
  -f                       Use fallback mode, uping will try to use next ping mode if selected by -p failed
  -c <count>               Number of pings to perform (default: 4)
  -i <ping-interval>       Interval between pings, e.g. -i 1s, -i 100ms (default: 1s)
  -t <host-timeout>        Timeout before probing one host terminates, regardless of how many pings perfomed, e.g. -t 1s, -t 100ms (default: <count> * 1s)
  -w <workers>             Number of parallel workers to run (default: 4)
  --source-db              Load hosts using database configured by -C <config-file>
  --source-api             Load hosts using external API configured by -C <config-file>
  --source-file <file-in>  Load hosts from file <file-in>
  --out-db                 Save tests results database configured by -C <config-file>
  --out-api                Save tests results using external API configured by -C <config-file>
  --out-file <file-out>    Save tests results to file <file-out>
`

//

const version = "0.3.5"

func loadHosts(hostsLoaders *[]schema.HostsLoader, hosts *schema.Hosts) {
	hosts.Reset()
	for _, hostsLoader := range *hostsLoaders {
		if err := hosts.Add(hostsLoader); err != nil {
			log.Fatal(err)
		}
	}
}

func pushJobs(jobs chan schema.Host, hosts *schema.Hosts) {
	for _, host := range hosts.Get() {
		jobs <- host
	}
}

func main() {
	var Hosts schema.Hosts

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], version)
	//fmt.Println(arguments)

	appConfig := schema.GeneralConfig{}
	hostsLoaders, resultsSavers, cleaners := configParser(arguments, &appConfig)

	// Load list of hosts
	Hosts.Init(appConfig.Probe.DefaultPort)
	loadHosts(&hostsLoaders, &Hosts)
	if len(Hosts.Get()) == 0 {
		log.Fatalln("No hosts to test.")
	}

	// Create workers pool
	jobs := make(chan schema.Host, appConfig.Workers)
	appConfig.Results = make(chan schema.ProbeResult, len(Hosts.Get()))

	var wgWorker sync.WaitGroup
	var wgWriter sync.WaitGroup

	wgWriter.Add(1)
	go worker.Saver(appConfig, resultsSavers, &wgWriter)

	for i := 0; i < appConfig.Workers; i++ {
		wgWorker.Add(1)
		go appConfig.Probe.Worker(i, appConfig, jobs, &wgWorker)
	}

	pushJobs(jobs, &Hosts)

	if appConfig.TestsInterval.Seconds() > 0.0 {
		ticker := time.NewTicker(appConfig.TestsInterval.Duration)

		for {
			select {
			case <-ticker.C:
				loadHosts(&hostsLoaders, &Hosts)
				pushJobs(jobs, &Hosts)
			}
		}
	}

	close(jobs)
	wgWorker.Wait()

	close(appConfig.Results)
	wgWriter.Wait()

	worker.Cleaner(appConfig, cleaners)
}
