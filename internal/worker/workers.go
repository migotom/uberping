package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/migotom/uberping/internal/schema"
	"github.com/migotom/uberping/internal/worker/netcat"
	goping "github.com/sparrc/go-ping"
)

// ResultsSaver saves PingResult.
type ResultsSaver func(schema.ProbeResult) error

func toMs(duration time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(duration.Nanoseconds())/1e6)
}

// Cleaner makes sure that all handles/sockets are closed before exiting app.
func Cleaner(config schema.GeneralConfig, cleaners []schema.HostsCleaner) {
	for _, cleaner := range cleaners {
		cleaner()
	}
}

// Saver worker iterates over config.Results tasks and saving them using ResultSavers functions.
func Saver(config schema.GeneralConfig, savers []ResultsSaver, wg *sync.WaitGroup) {
	defer wg.Done()

	var saversChannels []chan schema.ProbeResult
	var wgSavers sync.WaitGroup

	// assign results channel for each saver
	for _, resultSaver := range savers {
		ch := make(chan schema.ProbeResult, cap(config.Results))
		saversChannels = append(saversChannels, ch)

		// run each saver on separate gorutine
		wgSavers.Add(1)
		go func(results chan schema.ProbeResult, wgs *sync.WaitGroup, saver ResultsSaver) {
			for r := range results {
				if err := saver(r); err != nil {
					fmt.Println("error", err)
					log.Fatalln(err)
				}
			}
			wgs.Done()
		}(ch, &wgSavers, resultSaver)
	}

	// dispatch results between savers
	for result := range config.Results {
		for _, ch := range saversChannels {
			ch <- result
		}
	}

	// cleanup
	for i := range savers {
		close(saversChannels[i])
	}
	wgSavers.Wait()
}

// Netcat worer iterates over hosts tasks and try to establish to each of them connection to specified service.
func Netcat(id int, config schema.GeneralConfig, jobs <-chan schema.Host, wg *sync.WaitGroup) {
	defer wg.Done()

	for device := range jobs {
		var result schema.ProbeResult

		nc, err := netcat.NewNetcat(device.IP, device.Port)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		nc.OnFinish = func(stats *netcat.Statistics) {
			var line string
			if stats.ConnectionLoss < 100.0 {
				line = fmt.Sprintf("Connection to %s:%s succeded, time=%v!\n", stats.Addr, stats.Port, toMs(stats.Rtt))
			} else {
				line = fmt.Sprintf("Connection to %s:%s failed, %v!\n", stats.Addr, stats.Port, stats.ConnectionError)
			}
			result.Output = append(result.Output, line)
			result.Loss = stats.ConnectionLoss
			result.AvgTime = stats.Rtt.Seconds()
			result.Host = device
			config.Results <- result
		}

		nc.Timeout = config.Probe.Timeout.Duration

		nc.Run()
	}
}

// Pinger worker iterates over schema.Host tasks, running Ping command for each of them and push results into config.Results channel.
func Pinger(id int, config schema.GeneralConfig, jobs <-chan schema.Host, wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO move goping implementation to internal (needwd fix for 100ms dial connection timeout?)
	for device := range jobs {
		var result schema.ProbeResult

		pinger, err := goping.NewPinger(device.IP)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		pinger.OnRecv = func(pkt *goping.Packet) {

			line := fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, toMs(pkt.Rtt))

			if config.Verbose && !config.Grouped {
				fmt.Println(line)
			} else {
				result.Output = append(result.Output, line)
			}
		}

		pinger.OnFinish = func(stats *goping.Statistics) {
			var line string

			line += fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr)
			line += fmt.Sprintf("%d packets transmitted, %d packets received, %v packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			line += fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				toMs(stats.MinRtt), toMs(stats.AvgRtt), toMs(stats.MaxRtt), toMs(stats.StdDevRtt))

			result.Output = append(result.Output, line)
			result.Loss = stats.PacketLoss
			result.AvgTime = stats.AvgRtt.Seconds()
			result.Host = device
			config.Results <- result
		}

		pinger.SetPrivileged(config.Probe.Privileged)
		pinger.Interval = config.Probe.Interval.Duration
		pinger.Count = config.Probe.Count
		pinger.Timeout = config.Probe.Timeout.Duration

		pinger.Run()
	}
}
