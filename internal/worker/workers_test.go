package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/migotom/uberping/internal/schema"
)

func TestToMs(t *testing.T) {
	ms := toMs(time.Duration(1) * time.Second)
	if ms != "1000.000ms" {
		t.Errorf("toMs 1s conversion, got: %v, expected %v", ms, "1000.000ms")
	}
}

func TestSaver(t *testing.T) {
	var config schema.GeneralConfig
	config.Results = make(chan schema.PingResult, 1)

	var result string
	var s []ResultsSaver
	s = append(s, func(r schema.PingResult) error {
		result = "done"
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go Saver(1, config, s, &wg)

	config.Results <- schema.PingResult{Loss: 1}
	close(config.Results)
	wg.Wait()

	if result != "done" {
		t.Error("saver didn't save result")
	}

}

func TestPinger(t *testing.T) {
	var config schema.GeneralConfig

	// TODO consider, do we need to mockup this or test real behaviour like this?
	config.Ping.Privileged = false
	config.Ping.Count = 1
	config.Ping.Interval = time.Duration(1) * time.Second
	config.Ping.Timeout = time.Duration(1) * time.Second

	config.Results = make(chan schema.PingResult, 1)
	jobs := make(chan schema.Host, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go Pinger(1, config, jobs, &wg)

	jobs <- schema.Host{IP: "google.com"}
	result := <-config.Results
	close(jobs)
	close(config.Results)

	wg.Wait()

	if len(result.Output) < 2 {
		t.Errorf("missing pinger output, got: %v", result.Output)
	}
}
