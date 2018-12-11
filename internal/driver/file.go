package driver

import (
	"bufio"
	"fmt"
	"os"

	"github.com/migotom/uberping/internal/schema"
)

// FileLoadHosts loads list of hosts from file
func FileLoadHosts(hostParser schema.HostParser, filename string) ([]schema.Host, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []schema.Host
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip, port, err := hostParser(scanner.Text())
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, schema.Host{IP: ip, Port: port})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// FileSavePingResult save probe results to file.
func FileSavePingResult(result schema.ProbeResult, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range result.Output {
		fmt.Fprintln(file, line)
	}
	return nil
}
