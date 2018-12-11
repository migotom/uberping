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
	config.Results = make(chan schema.ProbeResult, 1)

	var result string
	var s []ResultsSaver
	s = append(s, func(r schema.ProbeResult) error {
		result = "done"
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go Saver(config, s, &wg)

	config.Results <- schema.ProbeResult{Loss: 1}
	close(config.Results)
	wg.Wait()

	if result != "done" {
		t.Error("saver didn't save result")
	}

}

func TestPinger(t *testing.T) {
	var config schema.GeneralConfig

	// TODO consider, do we need to mockup this or test real behaviour like this?
	config.Probe.Privileged = false
	config.Probe.Count = 1
	config.Probe.Interval.Duration = time.Duration(1) * time.Second
	config.Probe.Timeout.Duration = time.Duration(1) * time.Second

	config.Results = make(chan schema.ProbeResult, 1)
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
