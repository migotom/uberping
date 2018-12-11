package schema

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Host definition.
type Host struct {
	ID            int            `json:"id"`
	IP            string         `json:"ip"`
	Port          string         `json:"port"`
	InactiveSince sql.NullString `json:"inactive_since"`
}

// UnmarshalJSON is needed for unmarshal sq.NullString value used by SQL driver.
func (h *Host) UnmarshalJSON(data []byte) error {
	type Alias Host
	aux := &struct {
		InactiveSince string `json:"inactive_since"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	h.InactiveSince = sql.NullString{String: aux.InactiveSince, Valid: true}
	return nil
}

// HostParser validates input string as proper host and converts it to format accepted by probe.
type HostParser func(string) (string, string, error)

// HostsLoader returns list of hosts needed by probe workers, throws error in case failure of any validation.
type HostsLoader func(HostParser) ([]Host, error)

// HostsCleaner cleanups handlers, connections, open sockets, files etc. used by Loader/Saver/Parser.
type HostsCleaner func()

// Hosts defines list of hosts to probe.
type Hosts struct {
	hosts       []Host
	defaultPort int
}

func (h *Hosts) parseHost(host string) (string, string, error) {
	list := strings.Split(host, ":")
	if len(list) > 2 {
		return "", "", fmt.Errorf(fmt.Sprintf("Host invalid format: %s", host))
	}

	var port string
	if len(list) == 2 {
		port = list[1]
	} else {
		port = strconv.Itoa(h.defaultPort)
	}

	ipaddr, err := net.ResolveIPAddr("ip", list[0])
	if err == nil {
		return ipaddr.IP.String(), port, nil
	}

	IP, _, err := net.ParseCIDR(host)
	if err == nil {
		return IP.String(), port, nil
	}

	return "", "", fmt.Errorf(fmt.Sprintf("Can't resolve host: %s", host))
}

// Get list of hosts.
func (h *Hosts) Get() []Host {
	return h.hosts
}

// Init hosts list by setting up default port.
func (h *Hosts) Init(p int) {
	if p != 0 {
		h.defaultPort = p
	}
}

// Reset list of hosts.
func (h *Hosts) Reset() {
	h.hosts = nil
}

// Add hosts using HostsLoader function.
func (h *Hosts) Add(loader HostsLoader) error {
	hosts, err := loader(h.parseHost)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		h.hosts = append(h.hosts, host)
	}
	return nil
}
