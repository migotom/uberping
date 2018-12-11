package driver

import (
	"fmt"

	"github.com/migotom/uberping/internal/schema"
)

// StdoutPingResult saves probe results to STDOUT.
func StdoutPingResult(result schema.ProbeResult) error {
	for _, line := range result.Output {
		fmt.Println(line)
	}
	return nil
}
