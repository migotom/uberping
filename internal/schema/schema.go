package schema

import (
	"fmt"
	"net"
	"time"
)

// APIConfig defines external API settings.
type APIConfig struct {
	URL       string
	Name      string
	Secret    string
	AuthData  interface{}
	Endpoints APIEndpoints
}

// APIEndpoints defined extrnal API endpoints.
type APIEndpoints struct {
	Auth    string
	Devices string
	Device  string
}

// PingConfig sets up go-ping configuration.
type PingConfig struct {
	Privileged bool
	Interval   time.Duration
	Count      int
	Timeout    time.Duration
}

// PingResult keep result of go-ping operation.
type PingResult struct {
	Output  []string
	Loss    float64
	AvgTime float64
}

// GeneralConfig main application configuration.
type GeneralConfig struct {
	Verbose bool
	Grouped bool
	Workers int
	Results chan PingResult
	Ping    PingConfig
	API     APIConfig
}

// Host definition.
type Host struct {
	ID int    `json:"id"`
	IP string `json:"ip"`
}

// HostParser validates input string as proper host and converts it to format accepted by probe.
type HostParser func(string) (string, error)

// HostsLoader returns list of hosts needed by probe workers, throws error in case failure of any validation.
type HostsLoader func(HostParser) ([]Host, error)

// Hosts defines list of hosts to probe.
type Hosts struct {
	hosts []Host
}

func (h *Hosts) parseHost(host string) (string, error) {
	ipaddr, err := net.ResolveIPAddr("ip", host)
	if err == nil {
		return ipaddr.IP.String(), nil
	}

	IP, _, err := net.ParseCIDR(host)
	if err == nil {
		return IP.String(), nil
	}

	return "", fmt.Errorf(fmt.Sprintf("Can't resolve host: %s", host))
}

// Get list of hosts.
func (h *Hosts) Get() []Host {
	return h.hosts
}

// Set hosts using HostsLoader function.
func (h *Hosts) Set(loader HostsLoader) error {
	hosts, err := loader(h.parseHost)
	if err != nil {
		return err
	}

	h.hosts = hosts
	return nil
}
