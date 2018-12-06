package driver

import (
	"github.com/migotom/uberping/internal/schema"
)

// ArgvLoadHosts loads lists of hosts using standard argument list.
func ArgvLoadHosts(hostParser schema.HostParser, data []string) ([]schema.Host, error) {
	hosts := make([]schema.Host, len(data))
	for i, host := range data {
		ip, err := hostParser(host)
		if err != nil {
			return nil, err
		}

		hosts[i].IP = ip
	}
	return hosts, nil
}
