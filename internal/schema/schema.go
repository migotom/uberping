package schema

import (
	"time"
)

// PingConfig sets up go-ping configuration.
type PingConfig struct {
	Privileged bool
	Protocol   string
	Interval   Duration
	Count      int
	Timeout    Duration
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// PingResult keep result of go-ping operation.
type PingResult struct {
	Host    Host
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
	Results       chan PingResult
	Ping          PingConfig
	API           APIConfig
	DB            DBConfig
}
