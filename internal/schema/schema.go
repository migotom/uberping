package schema

import (
	"time"
)

// PingConfig sets up go-ping configuration.
type PingConfig struct {
	Privileged bool
	Interval   time.Duration
	Count      int
	Timeout    time.Duration
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
	TestsInterval time.Duration
	Workers       int
	Results       chan PingResult
	Ping          PingConfig
	API           APIConfig
	DB            DBConfig
}
