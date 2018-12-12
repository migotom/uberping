package schema

import (
	"sync"
	"time"
)

// Worker specifies worker type function.
type Worker func(id int, config GeneralConfig, jobs <-chan Host, wg *sync.WaitGroup)

// Duration is custom time.Duration implementation needed by TOML unmarshal.
type Duration struct {
	time.Duration
}

// UnmarshalText TOML config.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// ProbeConfig sets up go-ping configuration.
type ProbeConfig struct {
	Privileged  bool
	Mode        string
	Protocol    string
	Interval    Duration
	Count       int
	Timeout     Duration
	DefaultPort int `toml:"default_netcat_port"`
	Worker      Worker
}

// ProbeResult keep result of go-ping operation.
type ProbeResult struct {
	Host    Host
	Status  string
	Output  []string
	Loss    float64
	AvgTime float64
}

// GeneralConfig main application configuration.
type GeneralConfig struct {
	Verbose       bool
	Grouped       bool
	TestsInterval Duration `toml:"interval_between_tests"`
	Workers       int
	Results       chan ProbeResult
	Probe         ProbeConfig
	API           APIConfig
	DB            DBConfig
}
